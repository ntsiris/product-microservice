# Makefile

APP_NAME = product-microservice
VERSION = 1.0.0
DOCKER_IMAGE = $(APP_NAME):$(VERSION)
DOCKER_NETWORK = app_network
DOCKER_COMPOSE_FILE = docker/docker-compose.yml
DOCKER_FILE = docker/Dockerfile

.PHONY: all build run test docker-build docker-run docker-stop clean

all: test build

## Build the Go application binary
build:
	@echo "Building the application..."
	go build -o bin/$(APP_NAME) ./cmd

## Run the application locally
run:
	@echo "Running the application locally..."
	./bin/$(APP_NAME)

## Run tests with coverage
test:
	@echo "Running tests..."
	go test -cover ./...

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f $(DOCKER_FILE) .

## Run application using Docker Compose
docker-run:
	@echo "Running Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --build

## Stop and remove Docker Compose containers
docker-stop:
	@echo "Stopping Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

## Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f bin/$(APP_NAME)
