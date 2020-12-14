package redisHelpers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

type mockPool struct {
	conn *redigomock.Conn
}

func newMockPool() (mp *mockPool) {
	mp = &mockPool{
		conn: redigomock.NewConn(),
	}
	return
}

// Get returns a redis connection instance.
func (mp *mockPool) Get() (conn redis.Conn) {
	return mp.conn
}
