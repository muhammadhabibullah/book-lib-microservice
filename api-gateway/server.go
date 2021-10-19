package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	grpcClient "api-gateway/internal/grpc"
	httpHandler "api-gateway/internal/http"
	"api-gateway/internal/middleware"
	"api-gateway/pkg/proto"
)

const defaultPort = ":8000"

func init() {
	_ = godotenv.Load()
}

func main() {
	port := os.Getenv("HTTP_PORT")
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

	lendingGRPCClientConn, err := grpc.Dial(
		fmt.Sprintf("%s%s", os.Getenv("LENDING_SERVICE_HOST"), os.Getenv("LENDING_SERVICE_PORT")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	userServiceClient := proto.NewUserServiceClient(userGRPCClientConn)
	bookServiceClient := proto.NewBookServiceClient(bookGRPCClientConn)
	lendingServiceClient := proto.NewLendingServiceClient(lendingGRPCClientConn)

	userGRPCService := grpcClient.NewUserGRPCService(userServiceClient)
	bookGRPCService := grpcClient.NewBookGRPCService(bookServiceClient)
	lendingGRPCService := grpcClient.NewLendingGRPCService(lendingServiceClient)

	server := gin.Default()
	server.GET("/", httpHandler.GraphPlaygroundHandler())
	server.POST("/query", middleware.GinJWT(), httpHandler.GraphQLHandler(
		userGRPCService,
		bookGRPCService,
		lendingGRPCService,
	))

	if err = server.Run(defaultPort); err != nil {
		log.Fatal(err)
	}
}
