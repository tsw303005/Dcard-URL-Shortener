DOCKER_COMPOSE := $(or $(DOCKER_COMPOSE),$(DOCKER_COMPOSE),docker compose)

# clean
.PHONY: clean
clean:
	rm -rf bin/*

# build image
.PHONY: dc.image
dc.image: dc.build
	$(DOCKER_COMPOSE) build --force-rm image

.PHONY: dc.build
dc.build:
	$(DOCKER_COMPOSE) run --rm build

build:
	mkdir -p ./bin/app
	go build -o ./bin/app/cmd ./cmd/*.go

# lint
.PHONY: dc.pkg.lint
dc.pkg.lint:
	$(DOCKER_COMPOSE) run --rm lint make pkg.lint

.PHONY: dc.internal.lint
dc.internal.lint:
	$(DOCKER_COMPOSE) run --rm lint make internal.lint

.PHONY: dc.lint
dc.lint:
	$(DOCKER_COMPOSE) run --rm lint

pkg.lint:
	golangci-lint run ./pkg/...

internal.lint:
	golangci-lint run ./internal/...

lint:
	golangci-lint run ./...