package requestManager

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/stretchr/testify/assert"
)

func TestRequestManagerShouldTestRedisConnection(t *testing.T) {
	// PING successful
	pool := redisHelpers.NewMockPool()
	pingCmd := pool.Conn.Command("PING").Expect("PONG")

	ctx, cancel := context.WithCancel(context.Background())

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	_, err := Start(ctx, pool, logger)

	assert.Nil(t, err, "Shouldn't return an error when the PING succeeds")
	assert.True(t, pingCmd.Called, "Redis PING should be called when the request manager starts")

	cancel()

	// PING error
	pool = redisHelpers.NewMockPool()
	pingCmd = pool.Conn.Command("PING").ExpectError(fmt.Errorf("redis connection error"))

	ctx, cancel = context.WithCancel(context.Background())

	logBuffer = &bytes.Buffer{}
	logger = log.New(logBuffer, "", 0)

	_, err = Start(ctx, pool, logger)

	assert.Error(t, err, "Should return an error when the command fails.")

	cancel()
}

func TestRequestManagerShouldRespectContextCancellation(t *testing.T) {
	pool := redisHelpers.NewMockPool()
	pool.Conn.Command("PING").Expect("PONG")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	results, _ := Start(ctx, pool, logger)

	select {
	case _, open := <- results:
		assert.False(t, open, "Channel should be closed after the context is cancelled.")
	case <-time.After((config.PopTimeout+2)*time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time.")
	}

}
