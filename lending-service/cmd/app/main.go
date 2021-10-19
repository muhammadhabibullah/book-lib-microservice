package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"lending-service/internal/service"
	"lending-service/pkg/mongodb"
	"lending-service/pkg/proto"
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

	userGRPCClientConn, err := grpc.Dial(
		fmt.Sprintf("%s%s", os.Getenv("USER_SERVICE_HOST"), os.Getenv("USER_SERVICE_PORT")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	bookGRPCClientConn, err := grpc.Dial(
		fmt.Sprintf("%s%s", os.Getenv("BOOK_SERVICE_HOST"), os.Getenv("BOOK_SERVICE_PORT")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	userServiceClient := proto.NewUserServiceClient(userGRPCClientConn)
	bookServiceClient := proto.NewBookServiceClient(bookGRPCClientConn)

	lendingGRPCService := service.NewLendingGRPCService(userServiceClient, bookServiceClient)
	server := grpc.NewServer()
	proto.RegisterLendingServiceServer(server, lendingGRPCService)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
