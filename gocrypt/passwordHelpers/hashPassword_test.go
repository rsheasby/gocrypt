package passwordHelpers

import (
	"testing"

	"github.com/matryer/is"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordShouldSucceedForValidCost(t *testing.T) {
	is := is.New(t)

	pwd := []byte("password")
	cost := 10
	hash := HashPassword([]byte(pwd), cost)
	err := bcrypt.CompareHashAndPassword([]byte(hash), pwd)

	is.NoErr(err) // Hash and password should validate
}

// The reason for panicking with invalid cost is that it should be caught by the validation function in real operation.
func TestHashShouldPanicWithAboveMaxCost(t *testing.T) {
	is := is.New(t)

	defer func() {
		err := recover()
		if err == nil {
			is.Fail() // This should panic with an error if the cost is too high.
		}
	}()

	pwd := []byte("password")
	cost := 32
	_ = HashPassword([]byte(pwd), cost)
}
