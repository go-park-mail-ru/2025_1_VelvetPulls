SERVICE_NAME = app

COMPOSE_FILE = ./deploy/docker-compose.yml
ENV_FILE = .env

GO_CMD = go

COVERAGE_FILE = coverage.out

.PHONY: all
all: run

# Запуск приложения через Docker Compose
.PHONY: run
run:
	@docker-compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up --build $(SERVICE_NAME)

# Остановка контейнеров
.PHONY: stop
stop:
	@docker-compose -f $(COMPOSE_FILE) down

# Очистка файла покрытия
.PHONY: clean
clean:
	@rm -f $(COVERAGE_FILE)

# Тесты с покрытием
.PHONY: test
test:
	@$(GO_CMD) test ./tests/... \
		-coverpkg=./internal/usecase,./internal/repository,./internal/delivery/http,./internal/delivery/websocket,./pkg/middleware,./pkg/utils\
		-coverprofile=$(COVERAGE_FILE)
	@$(GO_CMD) tool cover -func=$(COVERAGE_FILE)
