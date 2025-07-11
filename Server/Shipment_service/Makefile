# Makefile

# Variables
DB_DRIVER ?= postgres
DB_USER ?= neondb_owner
DB_PASSWORD ?= npg_nU6TuylWD4kM
DB_HOST ?= ep-royal-night-a1q3qjxd-pooler.ap-southeast-1.aws.neon.tech
DB_PORT ?= 5432
DB_NAME ?= shipment_service
DB_SSLMODE ?= require
DB_DSN := $(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

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
