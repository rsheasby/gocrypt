package requestManager

import (
	"context"
	"log"
	"time"

	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

// Start starts the request manager, which pulls requests from redis, validates them, and puts them into the result channel.
func Start(ctx context.Context, pool redisHelpers.ConnGetter, logger *log.Logger) (results chan *protocol.Request, err error) {
	results = make(chan *protocol.Request, 1)

	if !config.Durable {
		// Test redis connection before going into the request loop
		conn := pool.Get()
		defer conn.Close()
		_, err = conn.Do("PING")
		if err != nil {
			return nil, err
		}
		_ = conn.Close()
	}

	go func() {
		for {
			if ctx.Err() != nil {
				close(results)
				return
			}
			req, err := redisHelpers.GetRequest(ctx, pool, logger)
			if err != nil {
				continue
			}
			err = validateRequest(req)
			if err != nil {
				logger.Printf("Invalid request received: %v", err)
				continue
			}
			expiryTime := time.Unix(0, req.ExpiryTimestamp)
			redisTime, _ := redisHelpers.GetRedisTime(pool)
			lateness := float64(redisTime.UnixNano()-expiryTime.UnixNano()) / 1000000000
			if lateness > 0 {
				logger.Printf(`Expired request received with response key "%s". It was %1.3f seconds late.`,
					req.ResponseKey, lateness)
				continue
			}
			results <- req
		}
	}()

	return
}
