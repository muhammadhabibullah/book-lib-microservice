package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"book-service/internal/service"
	"book-service/pkg/mongodb"
	"book-service/pkg/proto"
)

const defaultPort = ":8000"

func init() {
	_ = godotenv.Load()
}

func main() {
	mongodb.GetDatabase()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = defaultPort
	}

	bookService := service.NewBookGRPCService()
	server := grpc.NewServer()
	proto.RegisterBookServiceServer(server, bookService)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
