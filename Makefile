MAIN_FILE = cmd/main.go
GO_CMD = go
COVERAGE_FILE = coverage.out

.PHONY: all
all: run

.PHONY: run
run:
	@$(GO_CMD) run $(MAIN_FILE)

.PHONY: db_init
db_init:
	docker-compose --env-file .env build

.PHONY: clean
clean:
	rm $(COVERAGE_FILE)
	docker-compose down -v

.PHONY: test
test:
	@$(GO_CMD) test ./... -coverprofile=$(COVERAGE_FILE)

.PHONY: coverage
coverage:
	@$(GO_CMD) tool cover -func $(COVERAGE_FILE)
