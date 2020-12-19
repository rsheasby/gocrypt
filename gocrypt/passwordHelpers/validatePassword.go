package passwordHelpers

import (
	"golang.org/x/crypto/bcrypt"
)

// ValidatePassword takes a hash in unix "$2a" encoding and a password, and returns if the password is valid.
func ValidatePassword(password []byte, hash string) (isValid bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hash), password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
