package passwordHelpers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using the specified cost, and returns the hash in unix "$2a" encoding.
func HashPassword(password []byte, cost int) (hash string) {
	hashBytes, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		// Bcrypt only fails if something went very wrong, like OOM or a cost that's above the maximum.
		// Invalid cost should be caught by the validation, so if bcrypt fails, it's probably worth killing everything and
		// investigating.
		panic(err)
	}
	return string(hashBytes)
}
