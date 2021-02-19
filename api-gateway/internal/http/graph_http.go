package http_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"api-gateway/internal/domain/constant"
	"api-gateway/internal/graph"
	"api-gateway/internal/graph/generated"
	"api-gateway/internal/graph/model"
	grpcClient "api-gateway/internal/grpc"
)

func GraphPlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GraphQLHandler(
	userGRPCService *grpcClient.UserGRPCService,
	bookGRPCService *grpcClient.BookGRPCService,
	lendingGRPCService *grpcClient.LendingGRPCService,
) gin.HandlerFunc {
	h := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					UserGRPCService:    userGRPCService,
					BookGRPCService:    bookGRPCService,
					LendingGRPCService: lendingGRPCService,
				},
				Directives: generated.DirectiveRoot{
					IsAuthenticated: isAuthenticatedDirectiveConfig(),
					HasRole:         hasRoleDirectiveConfig(),
				},
			},
		),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func isAuthenticatedDirectiveConfig() func(
	ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		roleCtx := ctx.Value(constant.ClaimsGinCtxKey)
		if roleCtx == nil {
			return nil, errors.New("unauthorized")
		}
		return next(ctx)
	}
}

func hasRoleDirectiveConfig() func(
	ctx context.Context, obj interface{}, next graphql.Resolver, roles []*model.Role) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*model.Role) (interface{}, error) {
		roleCtx, ok := ctx.Value(constant.RoleGinCtxKey).(string)
		if !ok {
			return nil, fmt.Errorf("error parsing context: %s", constant.RoleGinCtxKey)
		}

		for _, role := range roles {
			if roleCtx == role.String() {
				return next(ctx)
			}
		}
		return nil, fmt.Errorf("unauthorized role: %s", roleCtx)
	}
}
