BUILD_DIR=./build

build: clean
	mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/authorizer

clean:
	rm -rf $(BUILD_DIR)

run:
	go run main.go

test: build	
	./test.fish $(BUILD_DIR)/authorizer

doc:
	@echo "Serving documentation...\n"
	@sleep 1 && echo "\nRead project documentation under http://localhost:6060/pkg/nuledger" &
	godoc -http=localhost:6060
