package main

import (
	"context"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/gocrypt/requestManager"
	"github.com/rsheasby/gocrypt/gocrypt/requestWorker"
)

func main() {
	// Read config from env vars
	config.ReadEnvironment()

	// Setup logger
	logOptions := log.Ldate | log.Ltime | log.Lmicroseconds
	if config.VerboseLogging {
		logOptions |= log.Llongfile
	}
	if config.UTCLogging {
		logOptions |= log.LUTC
	}
	logger := log.New(os.Stderr, "gocrypt:", logOptions)

	// Setup redis pool
	pool := &redis.Pool{
		MaxIdle:     config.Threads + 1,
		IdleTimeout: config.ConnectionTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				config.RedisHost,
				redis.DialUseTLS(config.RedisTLS),
				redis.DialUsername(config.RedisUsername),
				redis.DialPassword(config.RedisPassword),
				)
		},
	}

	// Open request manager. This exits the program if it's unable to connect to redis, unless Durable mode is enabled.
	requestChan, err := requestManager.Start(context.Background(), pool, logger)
	if err != nil {
		logger.Fatalf("Couldn't start up request manager: %v", err)
	}
	logger.Printf("gocrypt agent started, and Redis connection successfully opened to %s.", config.RedisHost)

	// Open request workers.
	requestWorker.StartMany(context.Background(), requestChan, pool, config.Threads, logger)

	// Let them do their work. If there's a fatal error, they will terminate the process.
	select {}
}
