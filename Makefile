
.PHONY: build run migrate-up migrate-down

build:
	go build -o bin/api.exe ./cmd/api

run:
	go run ./cmd/api

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down
