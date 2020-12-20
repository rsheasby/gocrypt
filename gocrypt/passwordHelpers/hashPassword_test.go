package passwordHelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordShouldSucceedForValidCost(t *testing.T) {
	pwd := []byte("password")
	cost := 10
	hash := HashPassword([]byte(pwd), cost)
	err := bcrypt.CompareHashAndPassword([]byte(hash), pwd)

	assert.NoError(t, err, "Generated hash should validate correctly")
}

// The reason for panicking with invalid cost is that it should be caught by the validation function in real operation.
func TestHashShouldPanicWithAboveMaxCost(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err, "Hashing with above max cost should panic with an error")
	}()

	pwd := []byte("password")
	cost := 32
	_ = HashPassword([]byte(pwd), cost)
}
