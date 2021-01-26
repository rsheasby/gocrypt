package remotePasswordHasher

const (
	// RequestQueueKey specifies the redis key that will be used for the request queue.
	RequestQueueKey = "gocrypt:RequestQueue"
	// ResponseKeyPrefix specifies the redis key prefix that will be used for response publishing.
	ResponseKeyPrefix = "gocrypt:Response:"
)
