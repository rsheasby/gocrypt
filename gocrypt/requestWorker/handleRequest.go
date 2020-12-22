package requestWorker

import (
	"log"

	"github.com/rsheasby/gocrypt/gocrypt/passwordHelpers"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

func handleRequest(request *protocol.Request, pool redisHelpers.ConnGetter, logger *log.Logger) {
	switch request.RequestType {
	case protocol.Request_HASHPASSWORD:
		handleHashRequest(request, pool, logger)
	case protocol.Request_VERIFYPASSWORD:
		handleValidateRequest(request, pool, logger)
	}
}

func handleHashRequest(req *protocol.Request, pool redisHelpers.ConnGetter, logger *log.Logger) {
	hash := passwordHelpers.HashPassword(req.Password, int(req.Cost))

	res := &protocol.Response{
		Hash: hash,
	}
	redisHelpers.PublishResponse(res, req.ResponseKey, pool, logger)
}

func handleValidateRequest(req *protocol.Request, pool redisHelpers.ConnGetter, logger *log.Logger) {
	isValid, err := passwordHelpers.ValidatePassword(req.Password, req.Hash)
	if err != nil {
		logger.Printf("Error when validating password: %v", err)
		return
	}

	res := &protocol.Response{
		IsValid: isValid,
	}
	redisHelpers.PublishResponse(res, req.ResponseKey, pool, logger)
}
