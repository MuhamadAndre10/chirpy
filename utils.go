package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, jwtSecret string, expiresIn time.Duration) (string, error) {

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy", // penerbit
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return ss, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	// 2. Parse token dengan klaim dan fungsi kunci
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing yang tidak valid: %v", t.Header["alg"])
		}

		return []byte(tokenSecret), nil

	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("gagal mengurai token: %w", err)
	}

	// 4. Periksa validitas token (tanda tangan dan klaim standar seperti exp)
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token tidak valid")
	}

	// dapatkan claims
	parsedClaims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("klaim token bukan tipe RegisteredClaims")
	}

	if parsedClaims.Subject == "" {
		return uuid.Nil, fmt.Errorf("klaim 'sub' kosong atau tidak ada")
	}

	userID, err := uuid.Parse(parsedClaims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("gagal mengurai subjek ('sub') sebagai UUID: %w", err)
	}

	return userID, nil

}

func GetBearerToken(header http.Header) (string, error) {

	// Get token value from header Authorization
	bearerToken := header.Get("Authorization")

	if bearerToken == "" {
		return "", fmt.Errorf("tidak ada value di header Authorization")
	}

	// bearerToken = Bearer TOKEN_STRING
	splitTokenStr := strings.Fields(bearerToken)

	if len(splitTokenStr) != 2 || splitTokenStr[0] != "Bearer" {
		return "", fmt.Errorf("format header Authorization tidak valid")
	}

	return splitTokenStr[1], nil

}
