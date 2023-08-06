package utils

import (
	"Momotok-Server/system"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(username string, uid int64) string {
	// create a token object
	token := jwt.New(jwt.SigningMethodHS256)
	// set claims of the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = uid
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	// generate toke string
	tokenString, _ := token.SignedString([]byte(system.ServerInfo.Server.SecretKey))

	return tokenString
}

func CheckToken(tokenString string) bool {
	if tokenString == "" {
		return false
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(system.ServerInfo.Server.SecretKey), nil //secret key
	})
	if token.Valid {
		//fmt.Println("token checked")
		return true
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return false
	} else {
		fmt.Println("Unable to handle token", err)
		return false
	}
}

func GetUsername(tokenString string) (string, error) {
	if CheckToken(tokenString) {
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 返回签名密钥
			return []byte(system.ServerInfo.Server.SecretKey), nil
		})
		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)
		return username, nil
	} else {
		err := errors.New("invalid token")
		return "", err
	}
}

func GetUID(tokenString string) (int64, error) {
	if CheckToken(tokenString) {
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 返回签名密钥
			return []byte(system.ServerInfo.Server.SecretKey), nil
		})
		claims := token.Claims.(jwt.MapClaims)
		uid := int64(claims["userid"].(float64))
		return uid, nil
	} else {
		err := errors.New("invalid token")
		return 0, err
	}
}

func GetUser(tokenString string) (int64, string, error) {
	if CheckToken(tokenString) {
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 返回签名密钥
			return []byte(system.ServerInfo.Server.SecretKey), nil
		})
		claims := token.Claims.(jwt.MapClaims)
		uid := int64(claims["userid"].(float64))
		username := claims["username"].(string)
		return uid, username, nil
	} else {
		err := errors.New("invalid token")
		return 0, "", err
	}
}
