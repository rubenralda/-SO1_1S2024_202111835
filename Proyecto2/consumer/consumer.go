package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Data struct {
	Name  string
	Album string
	Year  string
	Rank  string
}

func main() {
	fmt.Println("Version 4.1")
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
