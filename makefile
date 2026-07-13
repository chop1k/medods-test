# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Go & Tools
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
GO_COMPILER = go
GO_LINTER = go
GO_STATIC_ANALYZER = go
DOCKER_BIN = docker

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Application
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
APP_NAME = track-my-tasks
APP_VERSION = v0.0.0
APP_RELEASE = prod
APP_ROOT = .

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Application - Sources
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
SRC_DIR = internal
BIN_DIR = cmd
APP_SOURCES = $(shell find . -type f -name "*.go")

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Application - Testing
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
E2E_TESTS_DIR = $(APP_ROOT)/tests/e2e
E2E_BINARY_TARGET = $(APP_ROOT)/$(BUILD_DIR)/$(APP_VERSION)/e2e
UNIT_TEST_PACKAGES = ./internal/... ./cmd/...

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Application - Build
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
BUILD_DIR = build
BUILD_ENTRYPOINT = $(APP_ROOT)/$(BIN_DIR)/$(APP_NAME)/main.go
BUILD_TARGET = $(APP_ROOT)/$(BUILD_DIR)/$(APP_VERSION)/$(APP_NAME)

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Application - Deploy
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# compose.yml and Dockerfile live per-release under deploy/<release>/, but
# docker needs to see the whole repo (go.mod, internal/, cmd/, etc.) to
# build, so they get copied into the project root before running.
DEPLOY_DIR = deploy
COMPOSE_FILE = $(APP_ROOT)/compose.yml
DEPLOYED_DOCKERFILE = $(APP_ROOT)/Dockerfile

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# PHONY & such
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
PHONY_TARGETS = all
PHONY_TARGETS += all-tests all-build all-code
PHONY_TARGETS += run-binary run-compose run-docker
PHONY_TARGETS += run-tests-e2e run-tests-e2e-compose run-tests-unit
PHONY_TARGETS += tests-e2e tests-e2e-compose tests-unit
PHONY_TARGETS += application-binary application-compose application-docker-image
PHONY_TARGETS += deploy-prod deploy-test
PHONY_TARGETS += code-inspect code-style
PHONY_TARGETS += clean
PHONY_TARGETS += help

.PHONY: $(PHONY_TARGETS)
.DEFAULT_GOAL := all

all: all-tests all-build all-code

all-tests: tests-e2e tests-unit

all-build: application-binary application-compose application-docker-image

all-code: code-inspect code-style

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Run
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
run-binary: application-binary
	$(BUILD_TARGET) migrate && \
	$(BUILD_TARGET) seed && \
	$(BUILD_TARGET) serve

# Alias kept for whoever reaches for "compose" rather than "docker" - both
# bring the app up the same way, via the release's copied compose.yml.
run-compose: run-docker

run-docker: application-compose
	$(DOCKER_BIN) compose up

run-tests-e2e: $(E2E_BINARY_TARGET)
	$(E2E_BINARY_TARGET)

run-tests-e2e-compose: tests-e2e-compose
	$(DOCKER_BIN) compose up

run-tests-unit: tests-unit

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Tests
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
tests-e2e: $(E2E_BINARY_TARGET)

tests-e2e-compose: deploy-test

tests-unit:
	$(GO_COMPILER) test $(UNIT_TEST_PACKAGES)

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Build
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
application-binary: $(BUILD_TARGET)

application-compose: deploy-$(APP_RELEASE)

application-docker-image: application-compose
	$(DOCKER_BIN) build -f $(DEPLOYED_DOCKERFILE) -t $(APP_NAME):$(APP_VERSION) $(APP_ROOT)

$(BUILD_TARGET): $(APP_SOURCES)
	$(GO_COMPILER) build -o $(BUILD_TARGET) $(BUILD_ENTRYPOINT)

$(E2E_BINARY_TARGET): $(APP_SOURCES)
	$(GO_COMPILER) test -c -o $(E2E_BINARY_TARGET) $(E2E_TESTS_DIR)

# Copies the release's compose.yml and Dockerfile into the project root.
deploy-prod:
	cp $(DEPLOY_DIR)/prod/compose.yml $(COMPOSE_FILE)
	cp $(DEPLOY_DIR)/prod/Dockerfile $(DEPLOYED_DOCKERFILE)

deploy-test:
	cp $(DEPLOY_DIR)/test/compose.yml $(COMPOSE_FILE)
	cp $(DEPLOY_DIR)/test/Dockerfile $(DEPLOYED_DOCKERFILE)

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Code quality
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
code-inspect:
	$(GO_STATIC_ANALYZER) vet ./...

code-style:
	$(GO_LINTER) fmt ./...

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# Housekeeping
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(COMPOSE_FILE) $(DEPLOYED_DOCKERFILE)

help:
	@echo "Make scripts for track-my-tasks app:"
	@echo ""
	@echo "  all                     tests, build and code-quality checks (default)"
	@echo "  all-tests               build the e2e binary and run unit tests"
	@echo "  all-build               build the app binary, compose files and docker image"
	@echo "  all-code                run static checks and format code"
	@echo ""
	@echo "  run-binary              build the app and run migrate -> seed -> serve"
	@echo "  run-compose             alias of run-docker"
	@echo "  run-docker              copy deploy/$(APP_RELEASE) files to root and docker compose up"
	@echo ""
	@echo "  run-tests-unit          run unit tests"
	@echo "  run-tests-e2e           build and run the e2e test binary"
	@echo "  run-tests-e2e-compose   copy deploy/test files to root and docker compose up"
	@echo ""
	@echo "  code-inspect            run static analysis ($(GO_STATIC_ANALYZER) vet)"
	@echo "  code-style              format the codebase ($(GO_LINTER) fmt)"
	@echo ""
	@echo "  clean                   remove build artifacts and copied deploy files"
	@echo ""
	@echo "Configuration (override with VAR=value, e.g. make all GO_COMPILER=go1.22):"
	@echo "  GO_COMPILER             go binary to build/test with (default: $(GO_COMPILER))"
	@echo "  GO_LINTER               go binary used by code-style (default: $(GO_LINTER))"
	@echo "  GO_STATIC_ANALYZER      go binary used by code-inspect (default: $(GO_STATIC_ANALYZER))"
	@echo "  DOCKER_BIN              docker binary/path (default: $(DOCKER_BIN))"
	@echo "  APP_NAME                binary/image name (default: $(APP_NAME))"
	@echo "  APP_VERSION             build/image version (default: $(APP_VERSION))"
	@echo "  APP_RELEASE             prod or test - picks deploy/ config for application-compose (default: $(APP_RELEASE))"
	@echo "  APP_ROOT                project root, rarely needs changing (default: $(APP_ROOT))"
