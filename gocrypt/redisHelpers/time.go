package redisHelpers

import (
	"github.com/gomodule/redigo/redis"
)

// GetRedisTime returns the redis server's current timestamp
func GetRedisTime(pool ConnGetter) (timestamp int64) {
	conn := pool.Get()
	defer conn.Close()

	timestamps, err := redis.Int64s(conn.Do("TIME"))
	if err != nil {
		return -1
	}
	// This should never happen, since any error should actually have err set. Even so, doesn't hurt to put it here for safety.
	if len(timestamps) != 2 {
		return -1
	}

	return timestamps[0]
}