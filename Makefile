build:
	@go build -o bin/port-scanner

run:
	@./bin/port-scanner

test:
	@go test -v ./...