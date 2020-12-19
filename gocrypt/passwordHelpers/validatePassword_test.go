package passwordHelpers

import (
	"testing"

	"github.com/matryer/is"
)

func TestValidatePasswordShouldValidateCorrectPassword(t *testing.T) {
	is := is.New(t)

	password := []byte("password")
	hash := "$2y$10$z3QlTH2S0HFcX0bXY6B.8OQ.jj4mbdYPho4PnhEgM0qk2kbFNnw92"
	isValid, err := ValidatePassword(password, hash)

	is.NoErr(err)           // There shouldn't be an error for a valid password.
	is.Equal(isValid, true) // Should detect the password as valid.
}

func TestValidatePasswordShouldValidateIncorrectPassword(t *testing.T) {
	is := is.New(t)

	password := []byte("password1")
	hash := "$2y$10$z3QlTH2S0HFcX0bXY6B.8OQ.jj4mbdYPho4PnhEgM0qk2kbFNnw92"
	isValid, err := ValidatePassword(password, hash)

	is.NoErr(err)            // There shouldn't be an error for a invalid password with a valid hash.
	is.Equal(isValid, false) // Should detect the password as invalid.
}

func TestValidatePasswordShouldErrorWithInvalidHash(t *testing.T) {
	is := is.New(t)

	password := []byte("password")
	hash := "$2y$32$z3QlTH2S0HFcX0bXY6B.8OQ.jj4mbdYPho4PnhEgM0qk2kbFNnw92"
	isValid, err := ValidatePassword(password, hash)

	is.Equal(isValid, false) // Should report invalid if there is a hash error.
	if err == nil {
		is.Fail() // Should return an error when the hash is invalid.
	}
}
