DEV_REDIS_HOST = "127.0.0.1:6379"

DEV_BUILD_FLAGS =
RELEASE_BUILD_FLAGS = -ldflags "-s -w" -tags release
GLOBAL_BUILD_FLAGS = -o $(OUTPUT_FILE)

TEST_FLAGS = -v -race -count=1
TEST_TARGET = ./...

LINT_FLAGS = -E gosec -E gofmt --timeout 5m

TESTCOVERAGE_FILE = testcoverage.out
TESTCOVERAGE_FLAGS = $(RELEASE_BUILD_FLAGS) $(TEST_FLAGS) -cover -coverprofile $(TESTCOVERAGE_FILE) $(TEST_TARGET)

test:
	go test $(DEV_BUILD_FLAGS) $(TEST_FLAGS) $(TEST_TARGET)

test-release:
	go test $(RELEASE_BUILD_FLAGS) $(TEST_FLAGS) $(TEST_TARGET)

test-coverage:
	go test $(TESTCOVERAGE_FLAGS)

# To see test coverage:
# go get golang.org/x/tools/cmd/cover
show-coverage: test-coverage
	go tool cover -html=$(TESTCOVERAGE_FILE)

lint:
	golangci-lint run $(LINT_FLAGS)

fulltest: test test-release lint

