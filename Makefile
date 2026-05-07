include .env
export

.PHONY: up down migrate run dev

up:
	docker compose up -d

down:
	docker compose down

migrate:
	goose -dir database/migrations postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) sslmode=disable" up

run:
	go run cmd/server/main.go

dev: up migrate run
