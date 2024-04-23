package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	pb "server/proto"

	"google.golang.org/grpc"
)

var ctx = context.Background()

type server struct {
	pb.UnimplementedGetInfoServer
}

const (
	port = ":3001"
)

type Data struct {
	Name  string
	Album string
	Year  string
	Rank  string
}

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	fmt.Println("Recibí de cliente: ", in.GetName())
	data := Data{
		Name:  in.GetName(),
		Album: in.GetAlbum(),
		Year:  in.GetYear(),
		Rank:  in.GetRank(),
	}
	fmt.Println(data)
	// enviar el mensaje a kafka
	enviarMensaje(data)
	return &pb.ReplyInfo{Info: "Hola cliente, recibí el comentario"}, nil
}

func enviarMensaje(voto Data) {
	bytes, err := json.Marshal(voto)
	if err != nil {
		fmt.Println("Error al convertir a bytes:", err)
		return
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap:9092",
		"client.id":         "grcp_producer",
		"acks":              "all",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(bytes),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
