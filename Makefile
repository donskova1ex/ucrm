DEV_COMPOSE_ARGS=--env-file .env.local -f Docker-compose.yml
DEV_COMPOSE_ENV=docker compose $(DEV_COMPOSE_ARGS)
DEV_COMPOSE=docker compose $(DEV_COMPOSE_ARGS)


dev-build:
	$(DEV_COMPOSE) build
dev-up: dev-build
	$(DEV_COMPOSE) --env-file .env.local up -d