package remotePasswordHasher

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt"
	"golang.org/x/crypto/bcrypt"
)

// RemotePasswordHasher performs
type RemotePasswordHasher struct {
	cost int
	pool *redis.Pool
}

func testPoolConnection(pool *redis.Pool) (err error) {
	if pool == nil {
		return fmt.Errorf("redis pool cannot be nil")
	}
	testConn := pool.Get()
	if testConn == nil {
		// It doesn't seem like this ever happens,
		// as redigo opts to return the error when doing the actual operation instead of when getting the connection.
		// Irregardless, doesn't hurt to check it just in-case redigo changes, or my understanding is incorrect.
		return fmt.Errorf("nil connection returned from redis pool")
	}
	result, err := redis.String(testConn.Do("PING"))
	if err != nil {
		return fmt.Errorf("error PINGing redis: %v", err)
	}
	if result != "PONG" {
		// Unsure how to test this in a simple way, so it'll have to do without any coverage for now.
		return fmt.Errorf(`unexpected response when PINGing redis - expected "PONG", received "%s"`, result)
	}

	return nil
}

// NewRemotePasswordHasher returns a PasswordHasher instance relying on a remote gocrypt agent to perform the
// hashing. This validates the connection and cost, and returns an error if there is a problem.
func NewRemotePasswordHasher(cost int, pool *redis.Pool) (ph gocrypt.PasswordHasher, err error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("cost of %d is invalid - cost must be between %d and %d", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	testAttempts := 5
	for i := 0; i < testAttempts; i++ {
		err = testPoolConnection(pool)
		if err == nil {
			// Happy path
			return &RemotePasswordHasher{cost: cost, pool: pool}, nil
		}
	}
	// Bad path
	return nil, err
}

func (r RemotePasswordHasher) HashPassword(password string) (hash string, err error) {
	panic("implement me")
}

func (r RemotePasswordHasher) ValidatePassword(password string, hash string) (isValid bool, err error) {
	panic("implement me")
}
