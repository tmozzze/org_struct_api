# Makefile
.SILENT:

.PHONY: run build test lint clean

include .env
export

# Variables
# DSN FOR LOCAL
DSN := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_INTERNAL_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)"

APP_NAME=org_struct_api

# Install goose
goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

# MAKE Commands
run:
	go run cmd/api/main.go

test:
	go test -v ./...

lint:
	golangci-lint run --timeout 5m

# Swagger docs gen
swagger-gen:
	swag init -g cmd/api/main.go -o docs

up:
	docker-compose up --build -d

down:
	docker-compose down

down-and-clean:
	docker-compose down -v

# Migrations

# Create new migration (make create-migration NAME=name)
create-migration:
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
	@echo "Created migration: $(NAME)"

# Apply migrations
migrate-up:
	@echo "Applying migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DSN) up

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DSN) down

# Show migration status
migrate-status:
	@echo "Migration status:"
	goose -dir $(MIGRATIONS_DIR) postgres $(DSN) status

# Rollback all migrations
migrate-reset:
	@echo "Rolling back migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DSN) reset


# Debugging
debug:
	@echo "Current directory: $(shell pwd)"
	@docker-compose config