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
	// ConnectionTimeout specifies the timeout for the redis connection. Must be longer than the PopTimeout.
	ConnectionTimeout = 60 * time.Second
	// PopTimeout specifies the connection timeout for the blocking queue pop. This could be arbitrarily long, but you
	// have to set a limit so I reckon 10 seconds is reasonable.
	PopTimeout = 10
	// PublishAttempts specifies the maximum amount of times that the response publish will be retried if something goes wrong. Some tests rely on this being at least 3, so it should always be 3 or more.
	PublishAttempts = 5
	// MinResponseKeyLength specifies the minimum length for the response key. 16 is a decent length to be relatively sure you won't have collisions, and is also the length of a UUID
	MinResponseKeyLength = 16
)

var (
	// RedisHost specifies the host and port for the redis server.
	RedisHost string
	// RedisTLS specifies if the redis connection should use TLS.
	RedisTLS bool
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

	_, RedisTLS = os.LookupEnv("REDIS_TLS")
	if !RedisTLS {
		log.Println("Warning: TLS not enabled. Remember to configure and use TLS for any production deployments!")
	}

	Threads = runtime.NumCPU()
}
