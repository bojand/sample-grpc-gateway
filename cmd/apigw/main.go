package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"

	pb "github.com/bojand/sample-grpc-gateway/proto/sampleapi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type server struct {
	PublicURL string
}

func main() {
	grpcAddr := os.Getenv("SAMPLESVC_ADDR")

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
	srv := &server{PublicURL: publicURL}

	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// router
	r := chi.NewRouter()
	r.Get("/", srv.index)
	r.Mount("/api", gwmux) // all /api/* requests get routed to grpc service

	server := &http.Server{Addr: ":" + httpPort, Handler: r}

	log.Println("Serving gRPC-Gateway on " + httpAddr)
	log.Fatalln(server.ListenAndServe())
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
	
	log.Println("serve index")

	if err := tmpl.Execute(w, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
