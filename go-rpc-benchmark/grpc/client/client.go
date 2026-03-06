package main

import (
	"context"
	pb "grpc-service/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.GetUser(context.Background(), &pb.UserRequest{Id: "1"})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}
