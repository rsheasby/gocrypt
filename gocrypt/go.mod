module github.com/rsheasby/gocrypt/gocrypt

go 1.15

require (
	github.com/davecgh/go-spew v1.1.0
	github.com/gomodule/redigo v1.8.3
	github.com/google/uuid v1.1.2
	github.com/joho/godotenv v1.3.0
	github.com/rafaeljusto/redigomock v2.4.0+incompatible
	github.com/rsheasby/gocrypt/protocol v0.0.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
	google.golang.org/protobuf v1.25.0
)

replace github.com/rsheasby/gocrypt/protocol v0.0.0 => ../protocol
