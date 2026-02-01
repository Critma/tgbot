ENV_PATH=./configs/.env
include $(ENV_PATH)
COMPOSE_PATH=deploy/docker-compose.yml
MIGRATIONS_PATH=cmd/migrate/migrations
CMD_PATH=cmd/bot/main.go
CMD=build/main
TEST_DIR=build/test
COV_PATH=$(TEST_DIR)/coverage.out
COV_HTML_PATH=$(TEST_DIR)/index.html

compose-up:
	docker-compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) up -d

# run --
build-release:
	go build -o $(CMD) $(CMD_PATH)

run: build-release
	./$(CMD)

run-debug:
	go run $(CMD_PATH) -debug

run-ping:
	go run cmd/migrate/main.go

# --

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

# test --
test:
	@mkdir -p $(TEST_DIR)
	@go test -v ./... -coverprofile=$(COV_PATH)
	@go tool cover -html $(COV_PATH) -o $(COV_HTML_PATH)
	@firefox $(COV_HTML_PATH) || google-chrome $(COV_HTML_PATH) || echo "\nOpen $(COV_HTML_PATH) to see coverage info..."

# --

install-deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 

clean :
	rm -rf $(CMD) $(TEST_DIR)

.PHONY: run clean