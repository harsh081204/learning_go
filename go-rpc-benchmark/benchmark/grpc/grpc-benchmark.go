package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "grpc-service/proto"

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

	requests := 10000

	start := time.Now()

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			client.GetUser(context.Background(), &pb.UserRequest{Id: "1"})
		}()
	}

	wg.Wait()

	duration := time.Since(start)

	fmt.Println("gRPC requests:", requests)
	fmt.Println("Time:", duration)
	fmt.Println("RPS:", float64(requests)/duration.Seconds())
}
