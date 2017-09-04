DOCKER=$(shell which docker)
DOCKER_COMPOSE=$(shell which docker-compose)
CONTAINER_NAME=go-anki

up:
	$(DOCKER_COMPOSE) up -d

build:
	$(DOCKER_COMPOSE) build

clean: stop rm

attach:
	@$(DOCKER) exec -it $(CONTAINER_NAME) bash

stop:
	$(DOCKER_COMPOSE) stop

rm:
	$(DOCKER_COMPOSE) rm -f

run:
	@$(DOCKER) exec -i $(CONTAINER_NAME) go run $(SRC)

run-test:
	@$(DOCKER) exec -i $(CONTAINER_NAME) go test ./$(shell dirname $(SRC))
