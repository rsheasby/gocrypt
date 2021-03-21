package passwordHelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordShouldValidateCorrectPassword(t *testing.T) {
	password := []byte("password")
	hash := "$2y$10$z3QlTH2S0HFcX0bXY6B.8OQ.jj4mbdYPho4PnhEgM0qk2kbFNnw92"
	isValid, err := ValidatePassword(password, hash)

	assert.NoError(t, err, "Should not return error when validating a valid hash")
	assert.True(t, isValid, "Should validate as true when the correct password is provided")
}

func TestValidatePasswordShouldValidateIncorrectPassword(t *testing.T) {
	password := []byte("password1")
	hash := "$2y$10$z3QlTH2S0HFcX0bXY6B.8OQ.jj4mbdYPho4PnhEgM0qk2kbFNnw92"
	isValid, err := ValidatePassword(password, hash)

	assert.NoError(t, err, "Should not return error when validating a valid hash")
	assert.False(t, isValid, "Should validate as false when the incorrect password is provided")
}

func TestValidatePasswordShouldErrorWithInvalidHash(t *testing.T) {
	password := []byte("password")
	hash := "$2y$32$z3QlTH2S0HFcX0bXY6B.8OQ.jj4mbdYPho4PnhEgM0qk2kbFNnw92"
	isValid, err := ValidatePassword(password, hash)

	assert.False(t, isValid, "Should validate as false when there's an invalid hash provided")
	assert.NotNil(t, err, "Should return an error when an invalid hash is provided")
}
