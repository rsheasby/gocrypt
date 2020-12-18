package requestManager

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

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
			results <- req
		}
	}()

	return
}