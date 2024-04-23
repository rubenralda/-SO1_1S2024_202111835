package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Name  string `bson:"name,omitempty"`
	Album string `bson:"album,omitempty"`
	Year  string `bson:"year,omitempty"`
	Rank  string `bson:"rank,omitempty"`
}
type Logs struct {
	Name  string
	Fecha time.Time
}

func main() {
	fmt.Println("Version 5.1")
	// Conexion mongo
	uri := "mongodb://admin:1234@34.66.138.214:27017/?authSource=admin"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Conexion MONGO realizada")
	coll := client.Database("proyecto2").Collection("logs")
	newLogs := Logs{Name: "Prueba1", Fecha: time.Now().Local()}
	result, err := coll.InsertOne(context.TODO(), newLogs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ingresado", result)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap:9092",
		"group.id":          "grcp_producer",
		"auto.offset.reset": "earliest"})

	if err != nil {
		panic(err)
	}

	err = consumer.SubscribeTopics([]string{"myTopic"}, nil)
	run := true
	for run {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			voto := Data{}

			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			err := json.Unmarshal((*msg).Value, &voto)
			if err != nil {
				fmt.Println("Error al convertir:", err)
				continue
			}
			fmt.Println("Convertido", voto)
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()
}
