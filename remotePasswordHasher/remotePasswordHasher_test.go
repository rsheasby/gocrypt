package remotePasswordHasher

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// The tests in this file are both unit tests and integration tests.
// An active redis server on localhost:6379 is necessary, and a gocrypt agent needs to be running.

func TestNewRemotePasswordHasherShouldReturnPasswordHasher(t *testing.T) {
	cost := 10
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	ph, err := NewRemotePasswordHasher(cost, pool)

	assert.NotNil(t, ph, "Returned PasswordHasher shouldn't be nil")
	assert.Nil(t, err, "No error should be returned when the PasswordHasher was successfully created")

	// Ensure it validates the cost
	ph, err = NewRemotePasswordHasher(bcrypt.MinCost-1, pool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when the specified cost is below the minimum")
	assert.NotNil(t, err, "An error should be returned when the specified cost is below the minimum")

	ph, err = NewRemotePasswordHasher(bcrypt.MaxCost+1, pool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when the specified cost is above the maximum")
	assert.NotNil(t, err, "An error should be returned when the specified cost is above the maximum")

	// Ensure it tests the pool connection
	ph, err = NewRemotePasswordHasher(cost, nil)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when a nil pool is provided")
	assert.NotNil(t, err, "An error should be returned when a nil pool is provided")

	nilConnPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return nil, fmt.Errorf("no connection for you")
		},
	}
	ph, err = NewRemotePasswordHasher(cost, nilConnPool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when a connection can't be established")
	assert.NotNil(t, err, "An error should be returned when a connection can't be established")

	badHostPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "")
		},
	}
	ph, err = NewRemotePasswordHasher(cost, badHostPool)
	assert.Nil(t, ph, "PasswordHasher shouldn't be returned when a connection can't be established")
	assert.NotNil(t, err, "An error should be returned when a connection can't be established")
}
