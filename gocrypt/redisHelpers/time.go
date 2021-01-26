package redisHelpers

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// GetRedisTime returns the redis server's current timestamp
func GetRedisTime(pool ConnGetter) (redisTime time.Time, err error) {
	conn := pool.Get()
	defer conn.Close()

	timestamps, err := redis.Int64s(conn.Do("TIME"))
	if err != nil {
		return time.Time{}, fmt.Errorf("couldn't receive timestamp from redis: %v", err)
	}

	// Should never happen, but may as well check for it just in case
	if len(timestamps) != 2 {
		return time.Time{}, fmt.Errorf("couldn't receive timestamp from redis - invalid response")
	}

	return time.Unix(timestamps[0], timestamps[1]), nil
}
