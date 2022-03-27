PATH := $(CURDIR)/bin:$(PATH)

DOCKER_COMPOSE := $(or $(DOCKER_COMPOSE),$(DOCKER_COMPOSE),docker compose)

INTERNAL := internal

####################################################################################################
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

####################################################################################################
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

####################################################################################################
# generate

define make-dc-generate-rules

.PHONY: dc.$1.generate

# generate individual module rule
dc.$1.generate:
	$(DOCKER_COMPOSE) run --rm generate make $1.generate

endef
$(foreach module,$(INTERNAL),$(eval $(call make-dc-generate-rules,$(module))))

.PHONY: dc.pkg.generate
dc.pkg.generate:
	$(DOCKER_COMPOSE) run --rm generate make pkg.generate

.PHONY: dc.generate
dc.generate:
	$(DOCKER_COMPOSE) run --rm generate

define make-generate-rules

$1.generate: bin/mockgen
	go generate ./$1/...

endef
$(foreach module,$(INTERNAL),$(eval $(call make-generate-rules,$(module))))

pkg.generate: bin/mockgen
	go generate ./pkg/...

generate: pkg.generate $(addsuffix .generate,$(INTERNAL))

bin/mockgen: go.mod
	go build -o $@ github.com/golang/mock/mockgen

####################################################################################################
# test

define make-dc-test-rules

.PHONY: dc.$1.test
dc.$1.test:
	$(DOCKER_COMPOSE) run --rm make $1.test

endef
$(foreach module,$(INTERNAL),$(eval $(call make-dc-test-rules,$(module))))

.PHONY: dc.pkg.test
dc.pkg.test:
	$(DOCKER_COMPOSE) run --rm test make pkg.test

.PHONY: dc.test
dc.test:
	$(DOCKER_COMPOSE) run --rm test

define make-test-rules

$1.test:
	go test -v -race ./$1/...

endef
$(foreach module,$(INTERNAL),$(eval $(call make-test-rules,$(module))))

pkg.test:
	go test -v -race ./pkg/...

test: pkg.test $(addsuffix .test,$(INTERNAL))
