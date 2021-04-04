BUILD_DIR=./build

build: clean
	mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/authorizer

clean:
	rm -rf $(BUILD_DIR)

run:
	go run main.go

test:	
	go test

doc:
	@echo "Starting documentation server (godoc)..."
	@echo "Read project documentation under http://localhost:6060/pkg/nuledger\n"
	godoc -http=localhost:6060
