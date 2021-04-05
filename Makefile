BUILD_DIR=./build

export GOPROXY=https://proxy.golang.org,https://goproxy.io,direct

# Building and running locally

build: clean
	mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/authorizer

run:
	go run main.go

clean:
	rm -rf $(BUILD_DIR)

# Building and running locally with docker

docker:
	docker build -t authorizer .

docker_run: docker
	docker run -i --rm authorizer

# Tests and documentation

generate:
	go generate ./...

test:
	go test `go list ./... | grep -v mocks`

test_server:
	go run github.com/smartystreets/goconvey -excludedDirs=testcases,mocks,build

doc:
	@echo "Starting documentation server (godoc)..."
	@echo "Read project documentation under http://localhost:6060/pkg/nuledger\n"
	godoc -http=localhost:6060
