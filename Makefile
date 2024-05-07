all: test

tidy:
	go mod tidy
	go mod vendor

test:
	go test ./...

t-services:
	go test ./.../services/...

run:
	go run cmd/main.go