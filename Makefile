MAIN_FILE = main.go
GO_CMD = go
COVERAGE_FILE = coverage.out

.PHONY: all
all: run

.PHONY: run
run:
	@$(GO_CMD) run $(MAIN_FILE)

.PHONY: clean
clean:
	rm $(COVERAGE_FILE)

.PHONY: test
test:
	@$(GO_CMD) test ./... -coverprofile=$(COVERAGE_FILE)

.PHONY: coverage
coverage:
	@$(GO_CMD) tool cover -func $(COVERAGE_FILE)
