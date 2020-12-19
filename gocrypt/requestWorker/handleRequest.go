package requestWorker

import (
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/passwordHelpers"
	"github.com/rsheasby/gocrypt/gocrypt/redisHelpers"
	"github.com/rsheasby/gocrypt/protocol"
)

func handleRequest(request *protocol.Request, pool *redis.Pool, logger *log.Logger) {
	switch request.RequestType {
	case protocol.Request_HASHPASSWORD:
		handleHashRequest(request, pool, logger)
	case protocol.Request_VERIFYPASSWORD:
		handleValidateRequest(request, pool, logger)
	}
}

func handleHashRequest(req *protocol.Request, pool *redis.Pool, logger *log.Logger) {
	hash := passwordHelpers.HashPassword(req.Password, int(req.Cost))

	res := &protocol.Response{
		Hash:    hash,
	}
	redisHelpers.PublishResponse(res, req.ResponseKey, pool, logger)
}

func handleValidateRequest(request *protocol.Request, pool *redis.Pool, logger *log.Logger) {
	logger.Println("Not implemented") // TODO: implement me
}
