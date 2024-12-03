# Load .env file
include .env
export $(shell sed 's/=.*//' .env)

# Variables
MIGRATIONS_DIR := ./db/migrations

# Commands
.PHONY: tidy build run help create migrate-up migrate-down status rollback sqlc-gen

tidy:
	@go mod tidy

build:
	@go build -o cmd/bin/main cmd/main.go

run:
	@./cmd/bin/main

help:
	@echo "Available commands:"
	@echo "  make create name=<migration_name>   Create a new migration file"
	@echo "  make migrate-up                    Apply all up migrations"
	@echo "  make migrate-down                  Apply all down migrations"
	@echo "  make status                        Show the migration status"
	@echo "  make rollback                      Rollback the last migration"

create:
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' parameter is required. Usage: make create name=<migration_name>"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" down

status:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" status

rollback:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" down 1

sqlc-gen:
	sqlc generate