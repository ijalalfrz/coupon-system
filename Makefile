#===================#
#== Env Variables ==#
#===================#
DOCKER_COMPOSE_FILE ?= docker-compose.dev.yml
WITH_TTY ?= -t


migrate-up: ## Run migrations against non test DB
	APP_ENV=development docker compose -f ${DOCKER_COMPOSE_FILE} run --rm ${WITH_TTY} migrate up

migrate-down: ## Rollback migrations against non test DB
	APP_ENV=development docker compose -f ${DOCKER_COMPOSE_FILE} run --rm ${WITH_TTY} migrate down 1

migrate-create: ## Create a DB migration. You need to pass the file name, e.g. `make migrate-create name=migration-name`
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm migrate create -ext sql -dir /migrations $(name)

start:
	@echo "========================="
	@echo "Restarting services..."
	@echo "========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d --build
	docker compose -f ${DOCKER_COMPOSE_FILE} ps


stop:
	@echo "========================="
	@echo "Stopping services..."
	@echo "========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} down
	docker compose -f ${DOCKER_COMPOSE_FILE} ps


setup-env:
	@echo "========================="
	@echo "Setting up environment..."
	@echo "========================="
	cp ./.env.example ./.env

setup: setup-env start migrate-up

clean-db:
	rm -rf ./postgres-data

clean: stop clean-db

tests-unit:
	go test -v -timeout 10s -count=1 ./... -coverprofile=coverage.out

tests-load-all:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm ${WITH_TTY} k6 run /scripts/run_all_tests.js