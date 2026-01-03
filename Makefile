ENV_PATH=./configs/.env
COMPOSE_PATH=build/package/docker-compose.yml
MIGRATIONS_PATH=cmd/migrate/migrations
include $(ENV_PATH)

compose-up:
	docker-compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) up

# run
run-debug:
	go run cmd/bot/main.go -debug

run-migrate:
	go run cmd/migrate/main.go

# migrate --
migrate-create:
	migrate create -ext=sql -dir=$(MIGRATIONS_PATH) -seq init

migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database $(POSTGRES_URL) -verbose up 1

migrate-last:
	migrate -path=$(MIGRATIONS_PATH) -database $(POSTGRES_URL) -verbose up

migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database $(POSTGRES_URL) -verbose down 1

migrate-reset:
	migrate -path=$(MIGRATIONS_PATH) -database $(POSTGRES_URL) -verbose down

# --

install-deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 