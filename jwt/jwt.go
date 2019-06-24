package jwt

import (
	"crypto/rsa"
	"fmt"
	"github.com/MizukiSonoko/LndHub-go/logger"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

var (
	log        = logger.NewLogger()
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

//ToDo make it injectable using env
const (
	publicKeyPath  = "./public.key"
	privateKeyPath = "./secret.key"
	userIdClaimKey = "userId"
)

func init() {
	{
		key, err := ioutil.ReadFile(publicKeyPath)
		if err != nil {
			panic(fmt.Sprintf("ReadFile returns err:%s\n", err.Error()))
		}
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(key)
		if err != nil {
			panic(fmt.Sprintf("ParseRSAPublicKeyFromPEM returns err:%s\n", err.Error()))
		}
	}
	{
		key, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			panic(fmt.Sprintf("ReadFile returns err:%s\n", err.Error()))
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(key)
		if err != nil {
			panic(fmt.Sprintf("ParseRSAPrivateKeyFromPEM returns err:%s\n", err.Error()))
		}
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
		log.Error("token is invalid")
		return "", false
	}
	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims[userIdClaimKey].(string)
	if !ok {
		log.Error("request not set userId")
		return "", false
	}
	return userId, true
}

func GenerateToken(userId string) string {
	claims := make(jwt.MapClaims)
	claims[userIdClaimKey] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	fmt.Printf("tokenString:%s, err:%s\n", tokenString, err)
	return tokenString
}
