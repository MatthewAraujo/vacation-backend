package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseMultipartForm(r *http.Request) (map[string]any, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}
	payload := make(map[string]any)
	for key, value := range r.MultipartForm.Value {
		payload[key] = value[0]
	}
	for key, value := range r.MultipartForm.File {
		payload[key] = value[0]
	}
	return payload, nil
}
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(payload)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})

}
func Logger(status int, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	timestamp := time.Now().Format("02/01/2006 15:04:05")
	method := r.Method
	route := r.URL.Path
	fmt.Printf("%s [%s] %s %s %d\n", ip, timestamp, method, route, status)
}
