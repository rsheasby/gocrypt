module github.com/rsheasby/gocrypt

go 1.15

replace github.com/rsheasby/gocrypt/protocol v0.0.0 => ./protocol

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gomodule/redigo v1.8.3
	github.com/matryer/is v1.4.0 // indirect
	github.com/rsheasby/gocrypt/protocol v0.0.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	google.golang.org/protobuf v1.25.0
)
