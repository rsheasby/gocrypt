package config

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

const (
	// RequestQueueKey specifies the redis key that will be used for the request queue.
	RequestQueueKey = "gocrypt:RequestQueue"
	// ResponseKeyPrefix specifies the redis key prefix that will be used for response publishing.
	ResponseKeyPrefix = "gocrypt:Response:"
	// ErrorRetryTime specifies how long to wait before retrying when there's a redis read error.
	ErrorRetryTime = 1 * time.Second
	// PopTimeout specifies the connection timeout for the blocking queue pop. This could be arbitrarily long, but you
	// have to set a limit so I reckon 10 seconds is reaonable.
	PopTimeout = 10
)

var (
	// RedisHost specifies the host and port for the redis server.
	RedisHost string
	// Threads specifies how many worker threads should be started.
	Threads int
)

// ReadEnvironment gets the environment variables and initialises the config variables
func ReadEnvironment() {
	err := godotenv.Load("gocrypt.env")
	if err != nil {
		log.Println("Failed to read gocrypt.env. Falling back to environment variables.")
	}

	RedisHost = os.Getenv("REDIS_HOST")
	if RedisHost == "" {
		log.Fatalln(`No Redis host specified. Environment variable "REDIS_HOST" should be set.`)
	}

	Threads = runtime.NumCPU()
}
