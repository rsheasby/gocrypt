package requestManager

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
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
	pool.Conn.Command("PING").ExpectError(fmt.Errorf("redis connection error"))

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
	case _, open := <-results:
		assert.False(t, open, "Channel should be closed after the context is cancelled.")
	case <-time.After((config.PopTimeout + 2) * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time.")
	}
}

func TestRequestManagerShouldReturnValidRequestsWhileLoggingErrors(t *testing.T) {
	pool := redisHelpers.NewMockPool()
	pool.Conn.Command("PING").Expect("PONG")
	pool.Conn.Command("TIME").ExpectSlice(
		time.Now().Unix(),
		int64(time.Now().Nanosecond()),
	)

	req := &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Cost:            10,
		ExpiryTimestamp: math.MaxInt64,
	}
	reqBytes, _ := proto.Marshal(req)

	// Confirm that it doesn't break when there's an error in one of the requests.
	hasReturnedError := false
	hasTimedout := false
	pool.Conn.Command("BRPOP", config.RequestQueueKey, config.PopTimeout).Handle(func(args []interface{}) (interface{}, error) {
		if !hasReturnedError {
			hasReturnedError = true
			return nil, fmt.Errorf("Random error")
		}
		if !hasTimedout {
			hasTimedout = true
			return nil, redis.ErrNil
		}
		return []interface{}{
				[]byte(config.RequestQueueKey),
				reqBytes,
			},
			nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	results, err := Start(ctx, pool, logger)

	assert.Nil(t, err, "No error should be returned when starting the request manager")

	select {
	case req2 := <-results:
		// We need to compare the string values, as it fails otherwise due to the internal state differences.
		assert.EqualValues(t, req.String(), req2.String(), "Received request should be equal to the submitted request.")
	case <-time.After((config.PopTimeout + 2) * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time.")
	}
	assert.NotZero(t, logBuffer.Len(), "There should be logs confirming the error.")
}

func TestRequestManagerShouldValidateRequests(t *testing.T) {
	pool := redisHelpers.NewMockPool()
	pool.Conn.Command("PING").Expect("PONG")

	req := &protocol.Request{
		ExpiryTimestamp: math.MaxInt64,
	}
	reqBytes, _ := proto.Marshal(req)

	// Confirm that it doesn't break when there's an error in one of the requests.
	pool.Conn.Command("BRPOP", config.RequestQueueKey, config.PopTimeout).ExpectSlice([]byte(config.RequestQueueKey), reqBytes)

	ctx, cancel := context.WithCancel(context.Background())

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	results, err := Start(ctx, pool, logger)

	assert.Nil(t, err, "No error should be returned when starting the request manager")

	// This is pretty hacky, but I need to give the request manager enough time to actually receive the redis request.
	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case _, ok := <-results:
		assert.False(t, ok, "No message should be published since the request is invalid.")
	case <-time.After((config.PopTimeout + 2) * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time.")
	}

	assert.NotZero(t, logBuffer.Len(), "Should log when a request fails validation.")
}

func TestRequestManagerShouldCheckExpiryTime(t *testing.T) {
	pool := redisHelpers.NewMockPool()
	pool.Conn.Command("PING").Expect("PONG")
	pool.Conn.Command("TIME").ExpectSlice(int64(123), int64(0))

	// Request should be 23 seconds late
	req := &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUV",
		Password:        []byte("abc"),
		Cost:            10,
		ExpiryTimestamp: 100,
	}
	reqBytes, _ := proto.Marshal(req)

	// Confirm that it doesn't break when there's an error in one of the requests.
	pool.Conn.Command("BRPOP", config.RequestQueueKey, config.PopTimeout).ExpectSlice([]byte(config.RequestQueueKey), reqBytes)

	ctx, cancel := context.WithCancel(context.Background())

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	results, err := Start(ctx, pool, logger)

	assert.Nil(t, err, "No error should be returned when starting the request manager")

	// This is pretty hacky, but I need to give the request manager enough time to actually receive the redis request.
	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case _, ok := <-results:
		assert.False(t, ok, "No message should be published since the request is invalid.")
	case <-time.After((config.PopTimeout + 2) * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time.")
	}

	assert.NotZero(t, logBuffer.Len(), "Should log when a request was received too late.")
}
