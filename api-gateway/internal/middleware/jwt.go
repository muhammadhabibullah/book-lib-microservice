package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"api-gateway/internal/domain/constant"
	"api-gateway/pkg/jwt"
)

func GinJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.New().ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				map[string]interface{}{
					"errors": err.Error(),
				})
			return
		}
		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				map[string]interface{}{
					"errors": "token invalid",
				})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		id, _ := claims["id"].(string)
		role, _ := claims["role"].(string)
		email, _ := claims["email"].(string)

		requestCtx := ctx.Request.Context()
		requestCtx = context.WithValue(requestCtx, constant.ClaimsGinCtxKey, claims)
		requestCtx = context.WithValue(requestCtx, constant.UserIDGinCtxKey, id)
		requestCtx = context.WithValue(requestCtx, constant.RoleGinCtxKey, role)
		requestCtx = context.WithValue(requestCtx, constant.EmailGinCtxKey, email)
		ctx.Request = ctx.Request.WithContext(requestCtx)

		ctx.Next()
	}
}
