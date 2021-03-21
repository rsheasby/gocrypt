package redisHelpers

import (
	"context"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/protocol"
	"google.golang.org/protobuf/proto"
)

// GetRequest retrieves a hash request from redis. If no requests are currently in the queue, it blocks until one is available.
func GetRequest(ctx context.Context, pool ConnGetter, logger *log.Logger) (request *protocol.Request, err error) {
	// Continuously pop a request off one of the queues, retrying if IO timeout
	conn := pool.Get()
	defer conn.Close()

	for {
		if ctx.Err() != nil {
			return
		}
		result, err := redis.ByteSlices(conn.Do("BRPOP", config.RequestQueueKey, config.PopTimeout))
		if err == redis.ErrNil {
			continue
		}
		if err != nil {
			logger.Printf("Error receiving message from redis: %v", err)
			time.Sleep(config.ErrorRetryTime)
			return nil, err
		}
		// This should basically never happen. If there's no error, the response should always be 2 strings. Including this check just in case though.
		if len(result) != 2 {
			logger.Printf("Invalid response from Redis. Expected two strings but received %d.", len(result))
			continue
		}

		// Unmarshal the request received from redis
		request = &protocol.Request{}
		err = proto.Unmarshal(result[1], request)
		// It's unclear how bad a message has to be for proto.Unmarshall to fail, but I'm unable to make it happen, so this doesn't have any test coverage.
		if err != nil {
			logger.Printf("Failed to unmarshall message from redis: %s", err)
			continue
		}
		return request, err
	}
}
