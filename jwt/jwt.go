package jwt

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
)

var publicKey *rsa.PublicKey

//ToDo make it injectable using env
const (
	PUBLIC_KEY_PATH = "./sample.rsa.pub.pkcs8"
	userIdClaimKey  = "userId"
)

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

func GetUserIdFromToken(tokenStr string) (string, bool) {
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
