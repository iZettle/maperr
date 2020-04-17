.DEFAULT_GOAL := test

# -- Variables --

# The git tag for the version
TAG ?= $(shell git describe --exact-match --tags 2>/dev/null)

# -- High level targets --
# We only list these targets in the help. The other targets can still be used
# but it is generally better to call one of these.

## lint: run linter
.PHONY: lint
lint:
	$(call blue, "# running linter...")
	@golangci-lint run -c .golangci.yml

## test: run test suite for application
.PHONY: test
test:
	$(call blue, "# running tests...")
	@go test \
		-race ./...

## ci-test: run test suite for application
.PHONY: ci-test
ci-test:
	$(call blue, "# running tests...")
	@gotestsum --junitfile $(TEST_RESULTS_DIR)/unit-tests.xml -- \
		-race ./... \
		-coverprofile cover.out \
		-covermode atomic \
		-coverpkg ./...

.PHONY: help
help: ## Prints this help
ifneq ($(.DEFAULT_GOAL),)
	@echo "Default target: \033[36m$(.DEFAULT_GOAL)\033[0m"
endif
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# -- Low level targets --
# These targets are more low level and not included in the help. You can call
# them directly but generally you would use the higher level target.

.PHONY: clean
clean:
	@go clean -x