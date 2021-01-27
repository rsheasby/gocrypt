package remotePasswordHasher

import (
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/rsheasby/gocrypt"
	"github.com/rsheasby/gocrypt/protocol"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

// RemotePasswordHasher performs
type RemotePasswordHasher struct {
	cost    int
	timeout time.Duration
	pool    *redis.Pool
}

func testPoolConnection(pool *redis.Pool) (err error) {
	if pool == nil {
		return fmt.Errorf("redis pool cannot be nil")
	}
	testConn := pool.Get()
	if testConn == nil {
		// It doesn't seem like this ever happens,
		// as redigo opts to return the error when doing the actual operation instead of when getting the connection.
		// Irregardless, doesn't hurt to check it just in-case redigo changes, or my understanding is incorrect.
		return fmt.Errorf("nil connection returned from redis pool")
	}
	result, err := redis.String(testConn.Do("PING"))
	if err != nil {
		return fmt.Errorf("error PINGing redis: %v", err)
	}
	if result != "PONG" {
		// Unsure how to test this in a simple way, so it'll have to do without any coverage for now.
		return fmt.Errorf(`unexpected response when PINGing redis - expected "PONG", received "%s"`, result)
	}

	return nil
}

// New returns a PasswordHasher instance relying on a remote gocrypt agent to perform the
// hashing. This validates the connection and cost, and returns an error if there is a problem.
func New(cost int, timeout time.Duration, pool *redis.Pool) (ph gocrypt.PasswordHasher, err error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("cost of %d is invalid - cost must be between %d and %d", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	err = testPoolConnection(pool)
	if err != nil {
		return nil, err
	}

	return &RemotePasswordHasher{cost: cost, timeout: timeout, pool: pool}, nil
}

func generateResponseKey() (responseKey string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate a uuid response key: %v", err)
	}
	responseKey = fmt.Sprintf("%s-%d", id, time.Now().UnixNano())
	return
}

func encodePassword(password string) (encoded []byte) {
	shaBytes := sha512.Sum512([]byte(password))
	return shaBytes[:]
}

func (r RemotePasswordHasher) submitRequestAndGetResponse(req *protocol.Request) (res *protocol.Response, err error) {
	// Subscribe to hash res
	subConn := &redis.PubSubConn{
		Conn: r.pool.Get(),
	}
	defer subConn.Close()

	err = subConn.Subscribe(ResponseKeyPrefix + req.ResponseKey)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to res key: %v", err)
	}
	defer subConn.Unsubscribe()

	// Submit hash req
	redisConn := r.pool.Get()
	defer redisConn.Close()

	reqBytes, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall req: %v", err)
	}

	_, err = redisConn.Do("LPUSH", RequestQueueKey, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to submit hashing job: %v", err)
	}
	// Release connection early if successful. Double closing is safe. The deferred close is just for error conditions.
	redisConn.Close()

	// Receive hash res
	for {
		switch subResponse := subConn.ReceiveWithTimeout(r.timeout).(type) {
		case error:
			return nil, fmt.Errorf("failed to receive res from agent: %v", subResponse)
		case redis.Message:
			res = &protocol.Response{}
			err = proto.Unmarshal(subResponse.Data, res)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshall res from agent: %v", err)
			}

			return
		}
	}
}

func (r RemotePasswordHasher) getRedisTime() (redisTime time.Time, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	timestamps, err := redis.Int64s(conn.Do("TIME"))
	if err != nil {
		return time.Time{}, fmt.Errorf("couldn't receive timestamp from redis: %v", err)
	}

	// Should never happen, but may as well check for it just in case
	if len(timestamps) != 2 {
		return time.Time{}, fmt.Errorf("couldn't receive timestamp from redis - invalid response")
	}

	return time.Unix(timestamps[0], timestamps[1]), nil
}

func (r RemotePasswordHasher) HashPassword(password string) (hash string, err error) {
	responseKey, err := generateResponseKey()
	if err != nil {
		return "", fmt.Errorf("couldn't generate response key: %v", err)
	}

	redisTime, err := r.getRedisTime()
	if err != nil {
		return "", fmt.Errorf("couldn't get redis time: %v", err)
	}
	redisTime = redisTime.Add(r.timeout)

	req := &protocol.Request{
		RequestType:     protocol.Request_HASHPASSWORD,
		ResponseKey:     responseKey,
		Password:        encodePassword(password),
		Cost:            int32(r.cost),
		ExpiryTimestamp: redisTime.UnixNano(),
	}

	res, err := r.submitRequestAndGetResponse(req)
	if err != nil {
		return "", err
	}

	return res.Hash, nil
}

func (r RemotePasswordHasher) ValidatePassword(password string, hash string) (isValid bool, err error) {
	responseKey, err := generateResponseKey()
	if err != nil {
		return false, fmt.Errorf("couldn't generate response key: %v", err)
	}

	redisTime, err := r.getRedisTime()
	if err != nil {
		return false, fmt.Errorf("couldn't get redis time: %v", err)
	}
	redisTime = redisTime.Add(r.timeout)

	req := &protocol.Request{
		RequestType:     protocol.Request_VERIFYPASSWORD,
		ResponseKey:     responseKey,
		Password:        encodePassword(password),
		Hash:            hash,
		ExpiryTimestamp: redisTime.UnixNano(),
	}

	res, err := r.submitRequestAndGetResponse(req)
	if err != nil {
		return false, err
	}

	return res.IsValid, nil
}
