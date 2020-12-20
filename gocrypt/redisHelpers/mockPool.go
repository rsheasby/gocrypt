package redisHelpers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

type MockPool struct {
	Conn *redigomock.Conn
}

func NewMockPool() (mp *MockPool) {
	mp = &MockPool{
		Conn: redigomock.NewConn(),
	}
	return
}

// Get returns a redis connection instance.
func (mp *MockPool) Get() (conn redis.Conn) {
	return mp.Conn
}
