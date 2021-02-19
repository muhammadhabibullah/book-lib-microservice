package graph

import (
	grpcClient "api-gateway/internal/grpc"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	UserGRPCService    *grpcClient.UserGRPCService
	BookGRPCService    *grpcClient.BookGRPCService
	LendingGRPCService *grpcClient.LendingGRPCService
}
