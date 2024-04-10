package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	configs "github.com/MatthewAraujo/vacation-backend/config"
	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/MatthewAraujo/vacation-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handleFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the request
		tokenString := getTokenFromRequest(r)
		//validate the JWT
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("token is invalid")
			permissionDenied(w)
			return
		}
		// if is we need to fech the user from the DB
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userID, err := uuid.Parse(str)
		if err != nil {
			log.Printf("error parsing userID: %v", err)
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("error fetching user: %v", err)
			permissionDenied(w)
			return
		}
		// set context with the user

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)
		// call the next handler
		handleFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	// get the token from the request
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth == "" {
		return tokenAuth
	}
	return ""
}

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

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userID
}
