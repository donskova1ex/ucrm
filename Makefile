DEV_COMPOSE_ARGS=--env-file .env.local -f Docker-compose.yml
DEV_COMPOSE_ENV=docker compose $(DEV_COMPOSE_ARGS)
DEV_COMPOSE=docker compose $(DEV_COMPOSE_ARGS)

dev-build:
	$(DEV_COMPOSE) build
dev-up: dev-build
	$(DEV_COMPOSE) --env-file .env.local up -d

migrate-up:
	go run cmd/migrate/main.go -command=up

migrate-down:
	go run cmd/migrate/main.go -command=down

migrate-status:
	go run cmd/migrate/main.go -command=status

db-up:
	$(DEV_COMPOSE) up -d postgres

db-down:
	$(DEV_COMPOSE) down

db-logs:
	$(DEV_COMPOSE) logs postgres

run:
	go run cmd/api/main.go

build:
	go build -o bin/ucrm cmd/api/main.go