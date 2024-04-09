package auth

import (
	"time"

	configs "github.com/MatthewAraujo/vacation-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID string) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":  userID,
		"expires": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
