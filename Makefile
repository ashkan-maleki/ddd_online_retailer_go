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

ent-schema:
	go run -mod=mod entgo.io/ent/cmd/ent new --target ./internal/allocation/adaptors/ent Product


ent-new:
	@read -p 'Enter Schema Name: ' schema_name; echo; \
		ent new --target ./internal/ent/schema $$schema_name