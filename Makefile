all: test

tidy:
	go mod tidy
	go mod vendor

test:
	go test ./...


run:
	go run cmd/main.go