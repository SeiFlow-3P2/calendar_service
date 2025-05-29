package interceptor

import (
    "context"

    "google.golang.org/grpc"
)

func AuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Заглушка: реализуйте аутентификацию по необходимости
        return handler(ctx, req)
    }
}