package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {

	password := "test123"

	t.Run("hashPassword", func(t *testing.T) {

		newHashPass, err := HashPassword(password)

		assert.Nil(t, err, "no error")

		ok := ComparePasswordHash(password, newHashPass)

		assert.True(t, ok, "no error")

	})

	t.Run("test create jwt", func(t *testing.T) {

		superSecretCode := "qwertyuiop123"
		duration := 5 * time.Minute
		userID := uuid.New()

		token, err := MakeJWT(userID, superSecretCode, duration)
		assert.NoError(t, err, "no error")

		uID, err := ValidateJWT(token, superSecretCode)
		assert.NoError(t, err, "no error")

		assert.EqualValues(t, userID, uID, "userid harus sama")

		t.Logf("token, %v\n", token)
		t.Logf("err, %v\n", err)

	})

}

func TestAuth(t *testing.T) {
	testCase := []struct {
		name          string
		header        http.Header
		expectedToken string
		expectedError string
	}{
		{
			name:          "Sukses - Header Valid",
			header:        http.Header{"Authorization": []string{"Bearer my_secure_token_123"}},
			expectedToken: "my_secure_token_123",
			expectedError: "", // Tidak ada error
		},
		{
			name:          "Gagal - Header Kosong",
			header:        http.Header{}, // Header kosong
			expectedToken: "",
			expectedError: "tidak ada value di header Authorization",
		},
		{
			name:          "Gagal - Header Authorization Kosong",
			header:        http.Header{"Authorization": []string{""}}, // Value kosong
			expectedToken: "",
			expectedError: "tidak ada value di header Authorization",
		},
		{
			name:          "Gagal - Format Tidak Ada Bearer",
			header:        http.Header{"Authorization": []string{"Token my_secure_token_123"}}, // Tanpa "Bearer"
			expectedToken: "",
			expectedError: "format header Authorization tidak valid",
		},
		{
			name:          "Gagal - Hanya Bearer",
			header:        http.Header{"Authorization": []string{"Bearer"}}, // Hanya "Bearer"
			expectedToken: "",
			expectedError: "format header Authorization tidak valid",
		},
		{
			name:          "Sukses - Ada Spasi Lebih",
			header:        http.Header{"Authorization": []string{" Bearer   my_token_with_spaces "}}, // Spasi ekstra
			expectedToken: "my_token_with_spaces",                                                    // strings.Fields akan menangani spasi ekstra
			expectedError: "",
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Untuk kasus '%s', diharapkan error '%s', tapi dapat '%v'", tt.name, tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Untuk kasus '%s', tidak diharapkan error, tapi dapat '%v'", tt.name, err)
				}
				// Periksa apakah token yang dikembalikan sesuai harapan
				if token != tt.expectedToken {
					t.Errorf("Untuk kasus '%s', diharapkan token '%s', tapi dapat '%s'", tt.name, tt.expectedToken, token)
				}
			}
		})

	}

}

func TestToken(t *testing.T) {
	t.Run("test refresh token", func(t *testing.T) {

		token, err := MakeRefreshToken()

		t.Logf("token %s\n", token)

		assert.NoError(t, err, "no error")

	})

}
