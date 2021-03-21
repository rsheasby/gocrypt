package redisHelpers

import "github.com/gomodule/redigo/redis"

type ConnGetter interface {
	// Get returns a redis connection instance.
	Get() redis.Conn
}
