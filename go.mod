module github.com/rsheasby/gocrypt

go 1.15

replace github.com/rsheasby/gocrypt/protocol v0.0.0 => ./protocol

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gomodule/redigo v1.8.3
	github.com/rsheasby/gocrypt/protocol v0.0.0
	google.golang.org/protobuf v1.25.0
)
