package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

const JWTSECRET = "CPIFJWT666"

type CPIFUserClaims struct {
	Uuid     string
	UserID   string
	UserType int
	jwt.StandardClaims
}

// CreateJwtToken Create JWT Token
func CreateJwtToken(uuid, userID string, userType int) (string, error) {
	// 2 hours
	var TokenExpireDuration = time.Second * time.Duration(3600) * 2

	glog.Info("duration:", TokenExpireDuration)
	//jwt secert
	jwtSecret := []byte(JWTSECRET)

	c := CPIFUserClaims{
		uuid,
		userID,
		userType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "CPIF-SYSTEM",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(jwtSecret)
}

// ParseToken parse jwt token
func ParseToken(tokenString string) (*CPIFUserClaims, error) {
	// jwt secret
	var jwtSecret = []byte(JWTSECRET)
	// parse token
	token, err := jwt.ParseWithClaims(tokenString, &CPIFUserClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})

	if err != nil {
		glog.Errorln("parse error")
		return nil, err
	}

	fmt.Println("claims:")
	fmt.Println(token.Claims.(*CPIFUserClaims))

	if claims, ok := token.Claims.(*CPIFUserClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
