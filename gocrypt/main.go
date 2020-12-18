package main

import (
	"context"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt/gocrypt/config"
	"github.com/rsheasby/gocrypt/gocrypt/requestManager"
)

func main() {
	// Read config from env vars
	config.ReadEnvironment()

	// Setup logger
	var logger *log.Logger
	if config.UTCLogging {
		logger = log.New(os.Stderr, "gocrypt:",
			log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.LUTC)
	} else {
		logger = log.New(os.Stderr, "gocrypt:",
			log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	}

	// Setup redis pool
	pool := &redis.Pool{
		MaxIdle:     config.Threads + 1,
		IdleTimeout: config.ConnectionTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.RedisHost)
		},
	}

	// Test redis connection
	testConn := pool.Get()
	if err := testConn.Err(); err != nil {
		log.Fatalf("Unable to open Redis connection: %v", err)
	}
	if _, err := testConn.Do("PING"); err != nil {
		log.Fatalf("Redis connection not properly established: %v", err)
	}
	_ = testConn.Close()
	log.Println("Redis Connection Established.")

	requestChan := requestManager.Start(context.Background(), pool, logger)
	for req := range requestChan{
		spew.Dump(req)
	}
}
