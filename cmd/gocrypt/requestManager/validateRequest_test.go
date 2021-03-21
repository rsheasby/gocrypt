package requestManager

import (
	"math"
	"testing"

	"github.com/rsheasby/gocrypt/protocol"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequestShouldCatchGenericErrors(t *testing.T) {
	// Invalid request type
	req := &protocol.Request{
		RequestType:     3,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Hash:            "abc",
		Cost:            4,
		ExpiryTimestamp: math.MaxInt64,
	}

	err := validateRequest(req)
	assert.NotNil(t, err, "Should return an error when an invalid request type is provided")

	// Response key too short
	req = &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCD",
		Password:        []byte("abc"),
		Hash:            "abc",
		Cost:            4,
		ExpiryTimestamp: math.MaxInt64,
	}

	err = validateRequest(req)
	assert.NotNil(t, err, "Should return an error when response key is too short")

	// Password field empty
	req = &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte(""),
		Hash:            "abc",
		Cost:            4,
		ExpiryTimestamp: math.MaxInt64,
	}

	err = validateRequest(req)
	assert.NotNil(t, err, "Should return an error when password field is empty")
}

func TestValidateRequestShouldCatchHashPasswordErrors(t *testing.T) {
	// Low cost
	req := &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Hash:            "abc",
		Cost:            3,
		ExpiryTimestamp: math.MaxInt64,
	}

	err := validateRequest(req)
	assert.NotNil(t, err, "Should return an error when low cost provided")

	// High cost
	req = &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Hash:            "abc",
		Cost:            32,
		ExpiryTimestamp: math.MaxInt64,
	}

	err = validateRequest(req)
	assert.NotNil(t, err, "Should return an error when high cost provided")
}

func TestValidateRequestShouldCatchVerifyPasswordErrors(t *testing.T) {
	// Empty hash
	req := &protocol.Request{
		RequestType:     protocol.Request_VERIFYPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abc"),
		Hash:            "",
		Cost:            4,
		ExpiryTimestamp: math.MaxInt64,
	}

	err := validateRequest(req)
	assert.NotNil(t, err, "Should return an error when empty hash provided")
}

func TestValidateRequestShouldNotErrorWithValidRequest(t *testing.T) {
	// Valid hash request
	req := &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abd"),
		Cost:            10,
		ExpiryTimestamp: math.MaxInt64,
	}

	err := validateRequest(req)
	assert.Nil(t, err, "Should not error with valid hash request")

	// Valid verify request
	req = &protocol.Request{
		RequestType:     protocol.Request_VERIFYPASSWORD,
		ResponseKey:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Password:        []byte("abd"),
		Hash:            "abc",
		ExpiryTimestamp: math.MaxInt64,
	}

	err = validateRequest(req)
	assert.Nil(t, err, "Should not error with valid verify request")
}
