package middleware

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("HelloBrads")

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Okky"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatal("Something when wrong on JWT : ", err.Error())
		return "", nil
	}

	return tokenString, nil
}
