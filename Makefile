CURRENT_DIR := $(shell pwd)
LINT_BIN=go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint
GOLINE_BIN=go run github.com/segmentio/golines

#------------------------------------------------------------------------------
# Linter
#------------------------------------------------------------------------------

.PHONY: linecheck
linecheck:
	$(GOLINE_BIN) -m 110 -w ./

.PHONY: lint
lint:
	$(LINT_BIN) run

.PHONY: lint-fix
lint-fix: linecheck
	$(LINT_BIN) run --fix

# lint for debug
.PHONY: lint-with-cache
lint-with-cache:
	GOLANGCI_LINT_CACHE=$(CURRENT_DIR)/.cache/golangci-lint $(LINT_BIN) run

#------------------------------------------------------------------------------
# App
#------------------------------------------------------------------------------

.PHONY: run-v1
run-v1:
	go run ./cmd/v1/ --app app1 --count 100

.PHONY: run-v2
run-v2:
	go run ./cmd/v2/ --app app1 --count 100
