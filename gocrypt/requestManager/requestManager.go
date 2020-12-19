package requestManager

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

// Start starts the request manager, which pulls requests from redis, validates them, and puts them into the result channel.
func Start(ctx context.Context, pool *redis.Pool, logger *log.Logger) (results chan *protocol.Request) {
	results = make(chan *protocol.Request)

	// Test redis connection before going into the request loop
	conn := pool.Get()
	_, err := conn.Do("PING")
	if err != nil {
		logger.Fatalf("Failed to open test connection from request manager: %v", err)
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
			results <- req
		}
	}()

	return
}
