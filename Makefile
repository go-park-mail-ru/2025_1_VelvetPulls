MAIN_FILE = main.go
GO_CMD = go

.PHONY: all
all: run

.PHONY: run
run:
	@$(GO_CMD) run $(MAIN_FILE)

.PHONY: clean
clean:
	@echo "No build artifacts to clean."

.PHONY: test
test:
	@$(GO_CMD) test ./...