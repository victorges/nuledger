
build:
	go build

run:
	go run main.go

test: build	
	./test.fish

doc:
	@echo "Serving documentation...\n"
	@sleep 1 && echo "\nRead project documentation under http://localhost:6060/pkg/nuledger" &
	godoc -http=localhost:6060
