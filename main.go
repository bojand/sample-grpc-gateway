package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	"google.golang.org/grpc"

	pb "github.com/bojand/sample-grpc-gateway/proto/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer

	PublicURL string
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "5000"
	}
	grpcAddr := "0.0.0.0:" + grpcPort

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "3000"
	}
	httpAddr := "0.0.0.0:" + httpPort

	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://" + httpAddr
	}
	publicURL = strings.TrimSuffix(publicURL, "/")

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{PublicURL: publicURL}
	pb.RegisterGreeterServer(s, srv)

	log.Println("Serving gRPC on " + grpcAddr)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.index)
	mux.Handle("/api/hello", gwmux)

	gwServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on " + httpAddr)
	log.Fatalln(gwServer.ListenAndServe())
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	var (
		init sync.Once
		tmpl *template.Template
		err  error
	)
	init.Do(func() {
		tmpl, err = template.New("index").Parse(indexTemplate)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
