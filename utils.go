package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SuccJsonResponse(w http.ResponseWriter, code int, payload any) error {

	data, _ := json.Marshal(payload)

	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)

	return nil

}

func ErrJsonResponse(w http.ResponseWriter, code int, msg string) error {

	return SuccJsonResponse(w, code, map[string]string{"error": msg})

}

func HashPassword(password string) (string, error) {

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPass), nil

}

func ComparePasswordHash(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
