# Makefile

# Variables
DB_DRIVER ?= postgres
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= ticket_db
DB_SSLMODE ?= disable
DB_DSN := "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

# Default target
all: sqlc

# SQLC
sqlc:
	@echo ">> generating sqlc code..."
	sqlc generate

# Goose Migrations
goose-create-migration:
	goose -dir db/migrations create $(name) sql

goose-up:
	@echo ">> applying migrations..."
	goose -dir db/migrations postgres "$(DB_DSN)" up

goose-down:
	@echo ">> rolling back last migration..."
	goose -dir db/migrations postgres "$(DB_DSN)" down

goose-status:
	@echo ">> migration status..."
	goose -dir db/migrations postgres "$(DB_DSN)" status

.PHONY: all sqlc goose-create-migration goose-up goose-down goose-status
