package requestManager

import (
	"context"
	"log"

	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

// Start starts the request manager, which pulls requests from redis, validates them, and puts them into the result channel.
func Start(ctx context.Context, pool redisHelpers.ConnGetter, logger *log.Logger) (results chan *protocol.Request, err error) {
	results = make(chan *protocol.Request, 1)

	// Test redis connection before going into the request loop
	conn := pool.Get()
	_, err = conn.Do("PING")
	if err != nil {
		logger.Printf("Failed to open test connection from request manager: %v", err)
		return nil, err
	}
	_ = conn.Close()

	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			req, err := redisHelpers.GetRequest(ctx, pool, logger)
			if err != nil {
				continue
			}
			if err := validateRequest(req); err != nil {
				logger.Printf("Invalid request received: %v", err)
				continue
			}
			if lateness := redisHelpers.GetRedisTime(pool) - req.ExpiryTimestamp; lateness > 0 {
				logger.Printf(`Expired request received with response key "%s". It was %d seconds late.`, req.ResponseKey, lateness)
				continue
			}
			results <- req
		}
	}()

	return
}
