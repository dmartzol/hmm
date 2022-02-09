POSTGRES_HOST := localhost
POSTGRES_PORT := 5432
DB_NAME := development
POSTGRES_USER := user-development
POSTGRES_PASSWORD := password
POSTGRESQL_URL := postgresql://$(POSTGRES_HOST):$(POSTGRES_PORT)/$(DB_NAME)?user=$(POSTGRES_USER)&password=$(POSTGRES_PASSWORD)&sslmode=disable
MIGRATIONS_PATH := migrations
MIGRATE_VERSION := v4.15.1

.PHONY: up down migrate.up migrate.down

up:
	go mod tidy
	go mod download
	docker-compose up --remove-orphans -d --build

down:
	docker compose -p hmm down

migrate.up:
	docker run -v /Users/dani/dev/go/linkedin/migrations:/migrations \
				--network host migrate/migrate:$(MIGRATE_VERSION) \
				-path="$(MIGRATIONS_PATH)" \
				-database="$(POSTGRESQL_URL)" \
				-verbose \
				up

migrate.down:
	docker run -v /Users/dani/dev/go/linkedin/migrations:/migrations \
				--network host migrate/migrate:$(MIGRATE_VERSION) \
				-path="$(MIGRATIONS_PATH)" \
				-database="$(POSTGRESQL_URL)" \
				-verbose \
				down 1
