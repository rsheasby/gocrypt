package remotePasswordHasher

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/localPasswordHasher"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// The tests in this file are both unit tests and integration tests.
// An active redis server on localhost:6379 using TLS is necessary, and a gocrypt agent needs to be running.
// It also relies on the localPasswordHasher which serves as a reference implementation, so if that's broken,
// these tests can't be relied on.

func TestNewRemotePasswordHasher(t *testing.T) {
	cost := 10
	timeout := time.Second * 10
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379", redis.DialUseTLS(true))
		},
	}
	ph, err := New(cost, timeout, pool)

	assert.NotNil(t, ph, "Returned PasswordHasher shouldn't be nil")
	assert.Nil(t, err, "No error should be returned when the PasswordHasher was successfully created")

	// Ensure it validates the cost
	ph, err = New(bcrypt.MinCost-1, timeout, pool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when the specified cost is below the minimum")
	assert.NotNil(t, err, "An error should be returned when the specified cost is below the minimum")

	ph, err = New(bcrypt.MaxCost+1, timeout, pool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when the specified cost is above the maximum")
	assert.NotNil(t, err, "An error should be returned when the specified cost is above the maximum")

	// Ensure it tests the pool connection properly
	ph, err = New(cost, timeout, nil)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when a nil pool is provided")
	assert.NotNil(t, err, "An error should be returned when a nil pool is provided")

	nilConnPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return nil, fmt.Errorf("no connection for you")
		},
	}
	ph, err = New(cost, timeout, nilConnPool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when a connection can't be established")
	assert.NotNil(t, err, "An error should be returned when a connection can't be established")

	badHostPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "")
		},
	}
	ph, err = New(cost, timeout, badHostPool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when a connection can't be established")
	assert.NotNil(t, err, "An error should be returned when a connection can't be established")
}

func TestHashPassword(t *testing.T) {
	cost := 10
	timeout := time.Second * 10
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379", redis.DialUseTLS(true))
		},
	}

	lph, _ := localPasswordHasher.New(cost)
	rph, _ := New(cost, timeout, pool)

	// Test that it hashes passwords correctly
	pwd := "abc"
	hash, err := rph.HashPassword("abc")
	assert.Nil(t, err, "No error should be returned from the remote password hasher.")

	isValid, err := lph.ValidatePassword(pwd, hash)
	assert.True(t, isValid, "Hash from remote password hasher should validate using the local hasher.")
	assert.Nil(t, err, "No error should be returned when validating the hash. Invalid hash is likely the cause.")
}

func TestValidatePassword(t *testing.T) {
	cost := 10
	timeout := time.Second * 10
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379", redis.DialUseTLS(true))
		},
	}

	lph, _ := localPasswordHasher.New(cost)
	rph, _ := New(cost, timeout, pool)

	password := "password123!"
	hash, _ := lph.HashPassword(password)

	isValid, err := rph.ValidatePassword(password, hash)
	assert.Nil(t, err, "Validate password returned an error")
	assert.True(t, isValid, "Validate password didn't correctly validate")

	invalidPassword := "Password123!"
	isValid, err = rph.ValidatePassword(invalidPassword, hash)
	assert.Nil(t, err, "Validate password returned an error")
	assert.False(t, isValid, "Validate password didn't detect incorrect password")
}
