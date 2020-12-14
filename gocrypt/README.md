# gocrypt Agent
This is the service that performs all the password hashing on behalf of the client library.

## Setup
Firstly, you need a redis host running. You can do this locally like so:

```bash
docker pull redis:6
docker run -d -p 6379:6379 redis
```

This will download the latest Redis 6 Docker image, then run it with port 6379 exposed. You can verify the container is running using `docker ps`.

The connection to redis is specified using the `REDIS_HOST` environment variable. A default is provided in the makefile, so for convenience you can do `make run` to launch and connect to `localhost:6379`.

Note that this setup is **not** secure, so it should only be used for development.

## Communication
### Request
The gocrypt library and service communicate through Redis. The library will submit a request to either hash a new password or validate an existing hash using a `LPUSH` to the `gocrypt:RequestQueue` key. The gocrypt agent will `BRPOP` this key to receive requests. This essentially forms a FIFO queue of the password hash requests.

The requests are sent as a Protobuf message which is defined in the `protocol` directory.

### Response
The response is sent using Redis's Pub/Sub functionality. When the backend needs to submit a hash request, a large random key (like a UUID) is generated. This is submitted in the `response_key` parameter of the request. Before sending the request(to avoid race conditions), the library subscribes to the channel using that key in the following format: `gocrypt:Response:<response_key>`. When the agent is done with its hashing, it will publish the result using that key, which will be received by the backend.

As with the request, the response message is encoded in a Protobuf message as defined in the `protocol` directory.
