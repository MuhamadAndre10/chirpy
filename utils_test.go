package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHasingPassword(t *testing.T) {

	password := "test123"

	t.Run("hashPassword", func(t *testing.T) {

		newHashPass, err := HashPassword(password)

		assert.Nil(t, err, "no error")

		ok := ComparePasswordHash(password, newHashPass)

		assert.True(t, ok, "no error")

	})

}

func TestBcrypt(t *testing.T) {
	password := "testPassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		t.Errorf("Expected password to match, but got error: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongPassword"))
	if err == nil {
		t.Error("Expected password to not match, but it did")
	}
}
