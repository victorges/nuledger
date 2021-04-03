
build:
	go build

run:
	go run main.go

test: build	
	./test.fish

doc:
	@echo "Serving documentation..."
	@sleep 1 && echo "Read project documentation under http://localhost:6060/pkg/nuledger" &
	godoc -http=localhost:6060
