package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func loggingMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		fmt.Printf(
			"[%s] %s\n",
			time.Now().Format(time.RFC3339),
			info.FullMethod,
		)

		return handler(ctx, req)
	}
}
