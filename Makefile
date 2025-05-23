GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all
all: delimiter-AUDIT audit delimiter-UNIT-TESTS test delimiter-LINTERS lint delimiter-FINISH ## Runs multiple targets, audit, lint and test

.PHONY: audit
audit: ## Runs checks for security vulnerabilities on dependencies (including transient ones)
	go list -json -m all | nancy sleuth

.PHONY: build
build: ## Builds binary of library code as dis-redis
	go build ./...

.PHONY: convey
convey: ## Runs unit test suite and outputs results on http://127.0.0.1:8080/
	goconvey ./...

.PHONY: delimiter-%
delimiter-%:
	@echo '===================${GREEN} $* ${RESET}==================='

.PHONY: fmt
fmt: ## Run Go formatting on code
	go fmt ./...

.PHONY: lint
lint: ## Used in ci to run linters against Go code
	golangci-lint run ./...

.PHONY: lint-local
lint-local: ## Use locally to run linters against Go code
	golangci-lint run ./...

.PHONY: test
test: ## Runs unit tests including checks for race conditions and returns coverage
	go test -race -cover ./...

.PHONY: help
help: ## Show help page for list of make targets
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
