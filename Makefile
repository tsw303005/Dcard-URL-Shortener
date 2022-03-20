DOCKER_COMPOSE := $(or $(DOCKER_COMPOSE),$(DOCKER_COMPOSE),docker compose)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: dc.image
dc.image: dc.build
	$(DOCKER_COMPOSE) build --force-rm image

.PHONY: dc.build
dc.build:
	$(DOCKER_COMPOSE) run --rm build

build:
	mkdir -p ./bin/app
	go build -o ./bin/app/cmd ./cmd/*.go

