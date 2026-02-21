.PHONY: run build test lint clean

# Variables

APP_NAME=org_struct_api

# MAKE Commands

run:
	go run cmd/api/main.go

test:
	go test -v ./...

up:
	docker-compose up -d

down:
	docker-compose down