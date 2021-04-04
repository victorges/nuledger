BUILD_DIR=./build

export GOPROXY=https://proxy.golang.org,https://goproxy.io,direct

build: clean
	mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/authorizer

clean:
	rm -rf $(BUILD_DIR)

run:
	go run main.go

test:
	go test ./...

test_server:
	go run github.com/smartystreets/goconvey

generate:
	go generate ./...

doc:
	@echo "Starting documentation server (godoc)..."
	@echo "Read project documentation under http://localhost:6060/pkg/nuledger\n"
	godoc -http=localhost:6060
