BINARY := geo-service
MAIN := ./cmd/main.go
DOCKER_COMPOSE := docker-compose

.PHONY: all build run docker-build docker-up clean

all: build

build:
	@echo "Сборка бинарного файла $(BINARY)..."
	go build -o $(BINARY) $(MAIN)

run: build
	@echo "Запуск $(BINARY)..."
	./$(BINARY)

docker-build:
	@echo "Сборка Docker-образа..."
	$(DOCKER_COMPOSE) build

docker-up:
	@echo "Запуск контейнеров через docker-compose..."
	$(DOCKER_COMPOSE) up

clean:
	@echo "Очистка сгенерированных файлов..."
	rm -f $(BINARY)
