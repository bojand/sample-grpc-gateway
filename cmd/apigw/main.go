package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"

	pb "github.com/bojand/sample-grpc-gateway/proto/sampleapi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	grpcAddr := os.Getenv("SAMPLESVC_ADDR")

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "3000"
	}
	httpAddr := "0.0.0.0:" + httpPort

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

	fs := http.FileServer(http.Dir("./static"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", fs.ServeHTTP)
	mux.Handle("/api/hello", gwmux)
	mux.Handle("/api/reverse", gwmux)

	gwServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on " + httpAddr)
	log.Fatalln(gwServer.ListenAndServe())
}
