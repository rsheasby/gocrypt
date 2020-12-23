package localPasswordHasher

import (
	"crypto/sha512"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewLocalPasswordHasherShouldReturnPasswordHasher(t *testing.T) {
	cost := 10
	ph, err := NewLocalPasswordHasher(cost)

	assert.Nil(t, err, "NewLocalPasswordHasher shouldn't return any errors")
	assert.NotNil(t, ph, "NewLocalPasswordHelper shouldn't return nil")
}

func TestNewLocalPasswordHasherShouldPreventInvalidCosts(t *testing.T) {
	cost := bcrypt.MinCost - 1
	ph, err := NewLocalPasswordHasher(cost)

	assert.NotNil(t, err, "NewLocalPasswordHasher should return an error with a cost below the minimum.")
	assert.Nil(t, ph, "NewLocalPasswordHasher shouldn't return an instance with a cost below the minimum.")

	cost = bcrypt.MaxCost + 1
	ph, err = NewLocalPasswordHasher(cost)

	assert.NotNil(t, err, "NewLocalPasswordHasher should return an error with a cost above the maximum.")
	assert.Nil(t, ph, "NewLocalPasswordHasher shouldn't return an instance with a cost above the maximum.")
}

func TestPasswordHasherShouldHashPasswordsCorrectly(t *testing.T) {
	cost := 10
	ph, _ := NewLocalPasswordHasher(cost)

	pwd := "abc"
	hash, err := ph.HashPassword(pwd)
	assert.Nil(t, err, "Password hashing shouldn't return an err.")

	pwdSha := sha512.Sum512([]byte(pwd))
	err = bcrypt.CompareHashAndPassword([]byte(hash), pwdSha[:])
	assert.Nil(t, err, "Generated hash should validate with the provided password.")

	// Attempt to hash with empty password
	_, err = ph.HashPassword("")
	assert.NotNil(t, err, "Password hashing should return an error when an empty password is provided.")
}

func TestPasswordHasherShouldValidatePasswordsCorrectly(t *testing.T) {
	cost := 10
	ph, _ := NewLocalPasswordHasher(cost)

	pwd := "abc"
	pwdSha := sha512.Sum512([]byte(pwd))
	hash, _ := bcrypt.GenerateFromPassword(pwdSha[:], cost)

	isValid, err := ph.ValidatePassword(pwd, string(hash))
	assert.Nil(t, err, "Password validation shouldn't return an error with a valid hash.")
	assert.True(t, isValid, "Password should validate as correct.")

	wrongPwd := "ABC"
	isValid, err = ph.ValidatePassword(wrongPwd, string(hash))
	assert.Nil(t, err, "Password validation shouldn't return an error with a valid hash.")
	assert.False(t, isValid, "Password should validate as incorrect.")

	badHash := hash[:32]
	isValid, err = ph.ValidatePassword(pwd, string(badHash))
	assert.NotNil(t, err, "Password validation should return an error with an invalid hash.")
	assert.False(t, isValid, "Password should validate as incorrect with an invalid hash.")

	isValid, err = ph.ValidatePassword("", string(hash))
	assert.NotNil(t, err, "Password validation should return an error with an empty password provided.")
	assert.False(t, isValid, "Password should validate as incorrect with an empty password provided.")
}
