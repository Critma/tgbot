ENV_PATH=configs/.env
COMPOSE_PATH=build/package/docker-compose.yml

compose-up:
	docker-compose -f $(COMPOSE_PATH) --env-file $(ENV_PATH) up
