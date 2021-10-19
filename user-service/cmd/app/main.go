package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"user-service/internal/service"
	"user-service/pkg/mongodb"
	proto "user-service/pkg/proto"
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

	userService := service.NewUserGRPCService()
	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, userService)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
