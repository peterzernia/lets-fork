dc := docker-compose

build:
	$(dc) build
.PHONY: build

up:
	$(dc) up
.PHONY: up

lint:
	$(dc) run --rm app go vet
.PHONY: lint

test:
	$(dc) run --rm app go test -v ./...
.PHONY: test

app:
	$(dc) run --rm app go build
.PHONY: app

clean:
	$(dc) stop
	$(dc) rm -fv
.PHONY: clean
