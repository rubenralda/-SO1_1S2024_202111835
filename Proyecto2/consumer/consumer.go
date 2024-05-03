package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Logs struct {
	Name  string
	Fecha time.Time
}

func main() {
	fmt.Println("Version 6.7")
	// Conexion mongo
	uri := "mongodb://admin:1234@mongo-service.bds.svc.cluster.local:27017/?authSource=admin"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("proyecto2").Collection("logs")
	fmt.Println("Conexion MONGO realizada")

	//Conexion Redis
	client_redis := redis.NewClient(&redis.Options{
		Addr:     "redis-service.bds.svc.cluster.local:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	fmt.Println("Conexion REDIS realizada")

	// conexion kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092",
		"group.id":          "grcp_producer",
		"auto.offset.reset": "earliest"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Conexion KAFKA realizada")

	err = consumer.SubscribeTopics([]string{"myTopic"}, nil)
	run := true
	for run {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			userSession := client_redis.HGetAll(ctx, "votos").Val()
			valor, ok := userSession[string(msg.Value)]
			if !ok {
				valor = "0"
			}

			conteo, err := strconv.Atoi(valor)
			if err != nil {
				fmt.Println("Error al convertir")
				continue
			}

			err = client_redis.HSet(ctx, "votos", string(msg.Value), conteo+1).Err()
			if err != nil {
				fmt.Println("Error set redis:", err)
				continue
			}
			//Registrar logs de mensaje ingresado
			newLogs := Logs{Name: "Voto: " + string(msg.Value) + " Conteo: " + valor, Fecha: time.Now().Local()}
			result, err := coll.InsertOne(context.TODO(), newLogs)
			if err != nil {
				fmt.Println("Log no registrado: ", err)
				continue
			}
			fmt.Println("Ingresado:", result)
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()
}
