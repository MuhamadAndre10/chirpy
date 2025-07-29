package main

import (
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
