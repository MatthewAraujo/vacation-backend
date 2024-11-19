package xcsf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	configs "github.com/MatthewAraujo/vacation-backend/config"
	"github.com/MatthewAraujo/vacation-backend/utils"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/xcsf", GenerateXCSFToken).Methods(http.MethodGet)
}

func GenerateXCSFToken(w http.ResponseWriter, r *http.Request) {
	//do some logics to not allow anyone to generate this token
	csfToken, err := GenerateToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"csftoken": csfToken})

}

var secretKey = []byte(configs.Envs.XCSFToken)

func validateToken(token string) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid token")
	}

	payload := parts[0]
	signature := parts[1]

	expectedSignature := generateHMAC(payload, secretKey)

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return false, fmt.Errorf("invalid signature")
	}

	partsPayload := strings.Split(payload, "-")
	if len(partsPayload) < 2 {
		return false, fmt.Errorf("invalid payload")
	}

	timestamp := partsPayload[len(partsPayload)-1]
	createdAt, err := time.Parse("20060102150405", timestamp)
	if err != nil {
		return false, fmt.Errorf("error parsing timestamp: %v", err)
	}

	if time.Since(createdAt) > 5*time.Minute {
		return false, nil
	}

	return true, nil
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Xcsf")
	if token == "" {
		return ""
	}
	return token
}

func GenerateToken() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	randomID := hex.EncodeToString(randomBytes)

	timestamp := time.Now().Format("20060102150405")

	payload := fmt.Sprintf("%s-%s", randomID, timestamp)

	signature := generateHMAC(payload, secretKey)

	token := fmt.Sprintf("%s.%s", payload, signature)

	return token, nil
}

func generateHMAC(data string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func WithCSF(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromRequest(r)
		valid, err := validateToken(token)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if !valid {
			utils.PermissionDenied(w)
			return
		}

		handleFunc(w, r)
	}
}
