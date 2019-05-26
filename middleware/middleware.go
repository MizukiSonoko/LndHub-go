package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//ToDo make it injectable using env
const PUBLIC_KEY_PATH = ""
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

func WithJWT(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Header := r.Header.Get("Authorization")
		spd  := strings.SplitN(Header, " ", 2)
		if len(spd) != 2 || spd[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(spd[1], func(token *jwt.Token) (i interface{}, e error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method:%v",
					token.Header["alg"])
			}
			return publicKey, nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userId, ok := claims["userId"].(string)
		if !ok{
			log.Println("request not set userId")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userId)
		base.ServeHTTP(w, r.WithContext(ctx))
	})
}
