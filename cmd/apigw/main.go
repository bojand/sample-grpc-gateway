package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
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

	// Create a client connection to the gRPC server
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

	// grpc gateway mux
	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// static files
	fs := http.FileServer(http.Dir("./static"))

	// router
	r := chi.NewRouter()
	r.Get("/", fs.ServeHTTP)
	r.Mount("/api", gwmux) // all /api/* requests get routed to grpc service

	server := &http.Server{Addr: ":" + httpPort, Handler: r}

	log.Println("Serving gRPC-Gateway on " + httpAddr)
	log.Fatalln(server.ListenAndServe())
}
