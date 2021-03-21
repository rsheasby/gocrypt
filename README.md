# gocrypt
Opinionated Go library for scalable, secure password hashing
[![Go Reference](https://pkg.go.dev/badge/github.com/rsheasby/gocrypt.svg)](https://pkg.go.dev/github.com/rsheasby/gocrypt)

## How do I use it?
Firstly, you'll need a Redis server running. Assuming that's sorted, there's just 2 steps to start using gocrypt.

### Gocrypt agent
* `go get github.com/rsheasby/gocrypt/cmd/gocrypt`
* Create a `gocrypt.env` file according to the [example](https://github.com/rsheasby/gocrypt/blob/main/gocrypt/gocrypt.env)
* Run `gocrypt` in the same directory as the `gocrypt.env`

Alternatively, there is a docker image available [here](https://hub.docker.com/repository/docker/rsheasby/gocrypt). 
To configure the docker version, you can use environment variables.

### Gocrypt library
Once your agent is running, you can use the library by creating a redigo pool and creating a new 
RemotePasswordHasher instance:

```go
package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rsheasby/gocrypt"
	"github.com/rsheasby/gocrypt/remotePasswordHasher"
)

func main() {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	var ph gocrypt.PasswordHasher
	ph, _ = remotePasswordHasher.New(12, 30*time.Second, &pool)

	hash, _ := ph.HashPassword("hunter2")
	isValid, _ := ph.ValidatePassword("hunter2", hash)

	if isValid {
		log.Printf("Gocrypt is working!")
	}
}
```

Remember to add error handling :-)

## What does this do?
It provides an opinionated, simple, secure method of hashing passwords using separate hashing nodes that can be scaled independently of your backend. This keeps all your non-login requests responsive and fast since the hashing isn't hogging the CPU, and queues up all authentication requests to be executed in a scalable way, so that they can be distributed and dealt with as soon as more hashing power is available.

## What problem does this solve?
The goal of password hashing is to increase the computational cost of each password authentication, thereby slowing down a brute force attack in the event of a database breach or malicious employee with database access. However, this computational cost comes with a problem at scale. A well-tuned password hash should take at least half a second to complete(on one core), resulting in only a couple authentications per second. In some cases, these user authentications will come all at once, like when a mass email is sent out to your users, prompting them to login. This spike of requests can cause a significant increase in response times for all requests, and bring the whole application to a halt. Of course, you could just scale out your backend to spread out the load, but how many requests will timout in the time it takes to deploy new backend processes? When deploying new backend processes, the existing requests don't get redistributed to the new backends, so this is a reactionary measure. Also, even if you have enough backends, the long-running CPU task of hashing the passwords can increase the response time of other requests, as they wait for CPU time to become available. By queuing the hash requests to be run by dedicated runners, we avoid all these problems .

## What's under the hood?
gocrypt uses Redis for communication between the backend and the hashing nodes. SHA512 is used to hash the passwords before they are sent to redis. This provides basic obfuscation of the passwords in the queue, and allows arbitrary password lengths. Bcrypt is then used for the final password hashing. Bcrypt was chosen for its resistance to GPU acceleration, and the simple single parameter cost tuning. It's also taken a lot of cryptographic scrutiny, and it's held up pretty well so far.
