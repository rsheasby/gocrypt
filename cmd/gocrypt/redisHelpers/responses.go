package redisHelpers

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/protocol"
	"google.golang.org/protobuf/proto"
)

// PublishResponse publishes the provided response via redis, including automatic retry and responseKey concatenation with the prefix from the config package.
func PublishResponse(res *protocol.Response, responseKey string, pool ConnGetter, logger *log.Logger) {
	conn := pool.Get()
	defer conn.Close()

	resBytes, err := proto.Marshal(res)
	// This should never happen, but we'll check it for safety anyway
	if err != nil {
		logger.Printf(`Error publishing response "%s": Failed to marshall response: %v`, responseKey, err)
		return
	}

	for i := 1; i <= config.PublishAttempts; i++ {
		result, err := conn.Do("PUBLISH", config.ResponseKeyPrefix+responseKey, resBytes)
		if err != nil {
			logger.Printf(`Error publishing response "%s": Redis error when publishing response: %v`, responseKey, err)
			continue
		}

		receivedBy, err := redis.Int(result, nil)
		// Ditto
		if err != nil {
			logger.Printf(`Error publishing response "%s": Failed to interpret redis response: %v`, responseKey, err)
			return
		}
		if receivedBy == 0 {
			logger.Printf(`Error publishing response "%s": Published response wasn't received by any clients. Attempt %d of %d.`, responseKey, i, config.PublishAttempts)
			time.Sleep(config.ErrorRetryTime)
			continue
		}
		return
	}
	logger.Printf(`Error publishing response "%s": Unable to successfully publish response after %d attempt(s). Giving up.`, responseKey, config.PublishAttempts)
}
