package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtKey = []byte("042") //密钥，之后要用环境变量代替

type Claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims        //jwt包中的标准声明
}

func CreateToken(username string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Second) //token有效时间

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //token过期时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
