# Makefile for SARJAN

APP_NAME=sarjan
ENTRYPOINT=cmd/sarjan/main.go
DOCKER_IMAGE=sarjan:250829

.PHONY: all build run docker clean

all: build

build:
	go build -ldflags="-s -w" -o $(APP_NAME) $(ENTRYPOINT)

run: build
	./$(APP_NAME)

docker:
	docker build -t $(DOCKER_IMAGE) .

clean:
	rm -f $(APP_NAME)
