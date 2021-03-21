package requestWorker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

func TestRequestWorkerShouldHonorContextCancellation(t *testing.T) {
	pool := redisHelpers.NewMockPool()

	ctx, cancel := context.WithCancel(context.Background())

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	done := make(chan struct{})
	go func() {
		requestWorker(ctx, make(chan *protocol.Request), pool, logger)
		done <- struct{}{}
	}()

	cancel()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		assert.Fail(t, "Request worker didn't end in a reasonable time.")
	}
}

func TestRequestWorkerShouldProcessHashRequestsAndPublishTheResultCorrectly(t *testing.T) {
	t.Parallel()
	pool := redisHelpers.NewMockPool()

	hasErrored := false
	hasErroredWithNoRecipients := false

	doneChan := make(chan struct{})

	comm := pool.Conn.GenericCommand("PUBLISH").Handle(func(args []interface{}) (interface{}, error) {
		if !hasErrored {
			hasErrored = true
			return nil, fmt.Errorf("blah blah")
		}
		if !hasErroredWithNoRecipients {
			hasErroredWithNoRecipients = true
			return int64(0), nil
		}

		assert.Len(t, args, 2, "PUBLISH command should have 2 arguments")
		assert.Equal(t, config.ResponseKeyPrefix+"ABCDEFGHIJKLMNOPQRSTUVWXYZ", args[0], "PUBLISH key is incorrect")

		resBytes, ok := args[1].([]byte)
		assert.True(t, ok, "Response should be a byte array")

		res := &protocol.Response{}
		assert.Nil(t, proto.Unmarshal(resBytes, res), "Unmarshalling of response should succeed")

		assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(res.Hash), []byte("abc")), "Hash and password should validate")

		defer func() {
			doneChan <- struct{}{}
		}()
		return int64(1), nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	reqChan := make(chan *protocol.Request)

	StartMany(ctx, reqChan, pool, 1, logger)

	reqChan <- &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Cost:            int32(bcrypt.MinCost),
		ExpiryTimestamp: math.MaxInt64,
	}

	select {
	case <-doneChan:
		break
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time")
	}

	assert.True(t, comm.Called, "Request worker should publish the hash result.")
	assert.NotZero(t, logBuffer.Len(), "There should be some logs due to the simulated errors.")
}

func TestRequestWorkerShouldProcessVerifyValidRequestsAndPublishTheResultCorrectly(t *testing.T) {
	t.Parallel()
	pool := redisHelpers.NewMockPool()

	hasErrored := false
	hasErroredWithNoRecipients := false

	doneChan := make(chan struct{})

	comm := pool.Conn.GenericCommand("PUBLISH").Handle(func(args []interface{}) (interface{}, error) {
		if !hasErrored {
			hasErrored = true
			return nil, fmt.Errorf("blah blah")
		}
		if !hasErroredWithNoRecipients {
			hasErroredWithNoRecipients = true
			return int64(0), nil
		}

		assert.Len(t, args, 2, "PUBLISH command should have 2 arguments")
		assert.Equal(t, config.ResponseKeyPrefix+"ABCDEFGHIJKLMNOPQRSTUVWXYZ", args[0], "PUBLISH key is incorrect")

		resBytes, ok := args[1].([]byte)
		assert.True(t, ok, "Response should be a byte array")

		res := &protocol.Response{}
		assert.Nil(t, proto.Unmarshal(resBytes, res), "Unmarshalling of response should succeed")

		assert.True(t, res.IsValid, "Hash and password should validate")

		defer func() {
			doneChan <- struct{}{}
		}()
		return int64(1), nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	reqChan := make(chan *protocol.Request)

	StartMany(ctx, reqChan, pool, 1, logger)

	reqChan <- &protocol.Request{
		RequestType:     protocol.Request_VERIFYPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Hash:            "$2y$04$scoJ6DgfwqxqzQoTRdfvKOwQ1.aTPomv0rpoEub.FagPGAdvqW7Pa",
		ExpiryTimestamp: math.MaxInt64,
	}

	select {
	case <-doneChan:
		break
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time")
	}

	assert.True(t, comm.Called, "Request worker should publish the hash result.")
	assert.NotZero(t, logBuffer.Len(), "There should be some logs due to the simulated errors.")
}

func TestRequestWorkerShouldProcessVerifyInvalidRequestsAndPublishTheResultCorrectly(t *testing.T) {
	t.Parallel()
	pool := redisHelpers.NewMockPool()

	hasErrored := false
	hasErroredWithNoRecipients := false

	doneChan := make(chan struct{})

	comm := pool.Conn.GenericCommand("PUBLISH").Handle(func(args []interface{}) (interface{}, error) {
		if !hasErrored {
			hasErrored = true
			return nil, fmt.Errorf("blah blah")
		}
		if !hasErroredWithNoRecipients {
			hasErroredWithNoRecipients = true
			return int64(0), nil
		}

		assert.Len(t, args, 2, "PUBLISH command should have 2 arguments")
		assert.Equal(t, config.ResponseKeyPrefix+"ABCDEFGHIJKLMNOPQRSTUVWXYZ", args[0], "PUBLISH key is incorrect")

		resBytes, ok := args[1].([]byte)
		assert.True(t, ok, "Response should be a byte array")

		res := &protocol.Response{}
		assert.Nil(t, proto.Unmarshal(resBytes, res), "Unmarshalling of response should succeed")

		assert.False(t, res.IsValid, "Hash and password should validate")

		defer func() {
			doneChan <- struct{}{}
		}()
		return int64(1), nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	reqChan := make(chan *protocol.Request)

	StartMany(ctx, reqChan, pool, 1, logger)

	reqChan <- &protocol.Request{
		RequestType:     protocol.Request_VERIFYPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("Abc"),
		Hash:            "$2y$04$scoJ6DgfwqxqzQoTRdfvKOwQ1.aTPomv0rpoEub.FagPGAdvqW7Pa",
		ExpiryTimestamp: math.MaxInt64,
	}

	select {
	case <-doneChan:
		break
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Didn't receive a response within a reasonable time")
	}

	assert.True(t, comm.Called, "Request worker should publish the hash result.")
	assert.NotZero(t, logBuffer.Len(), "There should be some logs due to the simulated errors.")
}

func TestRequestWorkerShouldHandleVerifyRequestsWithInvalidHash(t *testing.T) {
	t.Parallel()
	pool := redisHelpers.NewMockPool()

	pool.Conn.GenericCommand("PUBLISH").Handle(func(args []interface{}) (interface{}, error) {
		assert.Fail(t, "Should not publish a response when an invalid hash is provided")
		return nil, nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	reqChan := make(chan *protocol.Request)

	StartMany(ctx, reqChan, pool, 1, logger)

	reqChan <- &protocol.Request{
		RequestType:     protocol.Request_VERIFYPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("Abc"),
		Hash:            "$2y",
		ExpiryTimestamp: math.MaxInt64,
	}

	time.Sleep(2 * time.Second)
	cancel()

	assert.NotZero(t, logBuffer.Len(), "There should be some logs due to the simulated errors.")
}

func TestRequestWorkerShouldAttemptToPublishTheCorrectAmountOfTimes(t *testing.T) {
	t.Parallel()
	pool := redisHelpers.NewMockPool()

	attempts := 0
	pool.Conn.GenericCommand("PUBLISH").Handle(func(args []interface{}) (interface{}, error) {
		attempts++

		return int64(0), nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "", 0)

	reqChan := make(chan *protocol.Request)
	doneChan := make(chan struct{})

	go func() {
		requestWorker(ctx, reqChan, pool, logger)
		doneChan <- struct{}{}
	}()

	reqChan <- &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Cost:            int32(bcrypt.MinCost),
		ExpiryTimestamp: math.MaxInt64,
	}

	cancel()

	select {
	case <-doneChan:
		assert.NotZero(t, logBuffer.Len(), "There should be some logs due to the simulated errors.")
		assert.Equal(t, config.PublishAttempts, attempts, "Didn't publish the correct amount of times.")
	case <-time.After((config.PublishAttempts + 1) * config.ErrorRetryTime):
		assert.Fail(t, "Didn't receive a response within a reasonable time.")
	}
}
