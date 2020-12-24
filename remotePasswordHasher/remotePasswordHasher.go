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
		return fmt.Errorf("nil connection returned from redis pool")
	}
	result, err := redis.String(testConn.Do("PING"))
	if err != nil {
		return fmt.Errorf("error PINGing redis: %v", err)
	}
	if result != "PONG" {
		return fmt.Errorf(`unexpected response when PINGing redis - expected "PONG", received "%s"`, result)
	}

	return nil
}

func NewRemotePasswordHasher(cost int, pool *redis.Pool) (ph gocrypt.PasswordHasher, err error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("cost of %d is invalid - cost must be between %d and %d", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	err = testPoolConnection(pool)
	if err != nil {
		return nil, err
	}

	return &RemotePasswordHasher{cost: cost, pool: pool}, nil
}

func (r RemotePasswordHasher) HashPassword(password string) (hash string, err error) {
	panic("implement me")
}

func (r RemotePasswordHasher) ValidatePassword(password string, hash string) (isValid bool, err error) {
	panic("implement me")
}
