package redisHelpers

import "github.com/gomodule/redigo/redis"

type connGetter interface {
	// Get returns a redis connection instance.
	Get() redis.Conn
}
