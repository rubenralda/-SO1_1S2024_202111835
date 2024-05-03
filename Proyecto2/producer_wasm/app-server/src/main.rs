use rocket::response::status::BadRequest;
use rocket::serde::json::{json, Value as JsonValue};
use rocket::serde::json::Json;
use rocket::config::SecretKey;
use rocket_cors::{AllowedOrigins, CorsOptions};
use std::time::Duration;
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::ClientConfig;


#[derive(rocket::serde::Deserialize)]
struct Data {
    name: String,
    album: String,
    year: String,
    rank: String,
}

async fn produce(data: &Data) -> Result<(), Box<dyn std::error::Error>> {
    // Configurar la dirección del broker y el nombre del tema Kafka
    let broker_address = "my-cluster-kafka-bootstrap:9092";
    let kafka_topic = "myTopic";
    
    // Configurar el cliente Kafka
    let producer: FutureProducer = ClientConfig::new()
        .set("bootstrap.servers", broker_address)
        .create()?;
    
    // Crear el mensaje a enviar
    let message_value = format!(
        r#"{{"name":"{}","album":"{}","year":"{}","rank":"{}"}}"#,
        data.name, data.album, data.year, data.rank
    );

    // Construir y enviar el mensaje
    let record = FutureRecord::to(kafka_topic)
        .key(&data.rank)
        .payload(&message_value);
    
    match producer.send(record, Duration::from_secs(0)).await {
        Ok(_) => println!("Mensaje enviado exitosamente"),
        Err((e, _)) => eprintln!("Error al enviar mensaje: {}", e),
    }

    Ok(())
}

#[rocket::post("/data", data = "<data>")]
fn receive_data(data: Json<Data>) -> Result<String, BadRequest<String>> {
    let received_data = data.into_inner();
    //let json_str = to_string(&received_data).map_err(|e| BadRequest(Some(e.to_string())))?;
    /* let response = JsonValue::from(json!({
        "message": format!("Received data: Sede: {}, Municipio: {}, Departamento: {}, Partido: {}", received_data.name, received_data.album, received_data.year, received_data.rank)
    })); */
    match produce(&received_data).await {
        Ok(_) => {
            let response = JsonValue::from(json!({
                "message": "Data received and sent to Kafka successfully"
            }));
            Ok(response.to_string())
        },
        Err(e) => {
            eprintln!("Error while producing message to Kafka: {}", e);
            Err(Status::InternalServerError)
        }
    }
}

#[rocket::main]
async fn main() {
    let secret_key = SecretKey::generate(); // Genera una nueva clave secreta

    // Configuración de opciones CORS
    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("failed to create CORS fairing");

    let config = rocket::Config {
        address: "0.0.0.0".parse().unwrap(),
        port: 8080,
        secret_key: secret_key.unwrap(), // Desempaqueta la clave secreta generada
        ..rocket::Config::default()
    };

    // Montar la aplicación Rocket con el middleware CORS
    rocket::custom(config)
        .attach(cors)
        .mount("/", rocket::routes![receive_data])
        .launch()
        .await
        .unwrap();
}