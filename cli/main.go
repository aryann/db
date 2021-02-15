package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/aryann/db/proto"
	"google.golang.org/grpc"
)

var (
	address = flag.String("address", ":8080", "The address of the server.")
)

func insert(ctx context.Context, client proto.DBClient, key string, payload string) {
	response, err := client.Insert(ctx, &proto.InsertRequest{
		Document: &proto.Document{
			Key:     key,
			Payload: payload,
		},
	})
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	log.Printf("Inserted %s at version %d.", response.GetVersion())
}

func main() {
	action := os.Args[1]

	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := proto.NewDBClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch action {
	case "insert":
		key := os.Args[2]
		payload := os.Args[3]
		insert(ctx, client, key, payload)
	default:
		log.Fatalf("Unknown command: %s", action)
	}
}
