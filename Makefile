dc := docker-compose

build:
	$(dc) build
.PHONY: build

up:
	$(dc) up
.PHONY: up

lint:
	$(dc) run --rm lets-fork go vet
.PHONY: lint

test:
	$(dc) run --rm lets-fork go test -v ./...
.PHONY: test

lets-fork:
	$(dc) run --rm lets-fork go build
.PHONY: lets-fork

clean:
	$(dc) stop
	$(dc) rm -fv
.PHONY: clean
