package main

import (
	"encoding/json"
	"net/http"
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
