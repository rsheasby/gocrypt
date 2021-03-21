module github.com/rsheasby/gocrypt/gocrypt

go 1.15

require (
	github.com/gomodule/redigo v1.8.3
	github.com/joho/godotenv v1.3.0
	github.com/rafaeljusto/redigomock v2.4.0+incompatible
	github.com/rsheasby/gocrypt v0.0.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620
	google.golang.org/protobuf v1.25.0
)

replace github.com/rsheasby/gocrypt => ../../
