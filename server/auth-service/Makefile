.PHONY: build run clean test fmt lint

# Сборка бинарника
build:
	cd ../dev-config/ && docker-compose up --build --force-recreate

run:
	cd ../dev-config/ && docker-compose up

# Удаление бинарника
clean:
	docker image prune

# Тестирование
test:
	go test -v ./...

# Форматирование
fmt:
	go fmt ./...

lint:
	golangci-lint run ./...