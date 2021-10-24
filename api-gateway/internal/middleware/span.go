package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func SpanContext(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		//span := tracing.StartSpanFromRequest(tracer, ginCtx.Request)
		//defer span.Finish()
		//
		//ctx := opentracing.ContextWithSpan(context.Background(), span)

		ginCtx.Next()
	}
}
