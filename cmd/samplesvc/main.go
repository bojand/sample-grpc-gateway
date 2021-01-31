package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/bojand/sample-grpc-gateway/proto/sampleapi"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("SayHello:", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *server) Reverse(ctx context.Context, in *pb.ReverseRequest) (*pb.ReverseResponse, error) {
	log.Println("Reverse:", in.GetInput(), "Upper:", in.GetUpper())
	reversed := reverse(in.GetInput())
	if in.GetUpper() {
		reversed = strings.ToUpper(reversed)
	}
	return &pb.ReverseResponse{Message: reversed}, nil
}

func main() {
	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = "5500"
	}
	grpcAddr := "0.0.0.0:" + grpcPort

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	log.Println("Serving gRPC on " + grpcAddr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func reverse(input string) string {
	// Get Unicode code points.
	n := 0
	rune := make([]rune, len(input))
	for _, r := range input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	// Convert back to UTF-8.
	return string(rune)
}
