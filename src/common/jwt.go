package common

import (
	"github.com/dgrijalva/jwt-go"
	"redisManger/src/models/LoginUserModel"
	"time"
)

var jwtKey = []byte("myKey")

type Claims struct {
	UserId   int64
	UserName string
	jwt.StandardClaims
}

func ReleaseToken(user LoginUserModel.LoginUser) (string, error) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId:   user.ID,
		UserName: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "admin",
			Subject:   "redisLogin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
