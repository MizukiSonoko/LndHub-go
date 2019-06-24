package interceptor

import (
	"context"
	"fmt"

	"github.com/MizukiSonoko/LndHub-go/controller"
	"github.com/MizukiSonoko/LndHub-go/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServiceAuthorize interface {
	Authorize(context.Context, string) error
}

func UnaryAuthenticateInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch info.FullMethod {
		case "/api.LndHubService/Login":
			return handler(ctx, req)
		default:
			newCtx, err := authentication(ctx)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
			return handler(newCtx, req)
		}
	}
}

func UnaryAuthorizationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		if srv, ok := info.Server.(ServiceAuthorize); ok {
			err = srv.Authorize(ctx, info.FullMethod)
		} else {
			return nil, fmt.Errorf("each service should implement an authorization")
		}
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return handler(ctx, req)
	}
}

func authentication(ctx context.Context) (context.Context, error) {
	fromMeta := func(ctx context.Context) (string, error) {
		data, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return "", fmt.Errorf("not found metadata")
		}
		vs := data["authorization"]
		if len(vs) == 0 {
			return "", fmt.Errorf("not found %s in metadata", "authorization")
		}
		return vs[0], nil
	}

	token, err := fromMeta(ctx)
	if err != nil {
		return nil, err
	}

	userId, ok := jwt.GetUserIdFromToken(token)
	if !ok {
		return nil, fmt.Errorf("GetUserIdFromToken not ok")
	}

	c := context.WithValue(ctx, controller.CtxUserIdKey, userId)
	fmt.Printf("res c:%v\n", c)
	return c, nil
}
