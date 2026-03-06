package main

import (
	"context"
	"log"
	"net"

	pb "grpc-service/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &server{})

	log.Println("gRPC running at :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
