package interceptor

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
)

//ToDo make it injectable using env
const (
	PUBLIC_KEY_PATH = ""
	userIdClaimKey  = "userId"
)

var publicKey *rsa.PublicKey

func init() {
	key, err := ioutil.ReadFile(PUBLIC_KEY_PATH)
	if err != nil {
		panic(fmt.Sprintf("ReadFile returns err:%s\n", err.Error()))
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		panic(fmt.Sprintf("ParseRSAPublicKeyFromPEM returns err:%s\n", err.Error()))
	}
}

type ServiceAuthorize interface {
	Authorize(context.Context, string) error
}

func UnaryAuthenticateInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := authentication(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return handler(newCtx, req)
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
			return "", fmt.Errorf("not found %s in metadata",  "authorization")
		}
		return vs[0], nil
	}

	token, err := fromMeta(ctx)
	if err != nil {
		return nil, err
	}

	userId, ok := getUserIdFromToken(token)
	if !ok {
		return nil, err
	}

	return context.WithValue(ctx, "userId", userId), nil
}

func getUserIdFromToken(tokenStr string) (string, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method:%v",
				token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return "", false
	}
	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims[userIdClaimKey].(string)
	if !ok {
		log.Println("request not set userId")
		return "", false
	}
	return userId, true
}

func GenerateToken(userId string) string {
	claims := make(jwt.MapClaims)
	claims[userIdClaimKey] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.Raw
}
