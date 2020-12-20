package requestWorker

import (
	"context"
	"log"

	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

// Start starts the specified amount of request workers to receive and process requests, then publish the results back to the client via redis.
func StartMany(ctx context.Context, reqChan chan *protocol.Request, pool redisHelpers.ConnGetter, count int, logger *log.Logger) {
	for i := 0; i < count; i++ {
		go requestWorker(ctx, reqChan, pool, logger)
	}
	logger.Printf("Started %d worker thread(s).", count)
}

func requestWorker(ctx context.Context, reqChan chan *protocol.Request, pool redisHelpers.ConnGetter, logger *log.Logger) {
	for {
		// This is duplicated so that a cancelled context takes priority over the request channel.
		if ctx.Err() != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		case req := <- reqChan:
			handleRequest(req, pool, logger)
		}
	}
}