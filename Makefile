LOCAL_BIN := $(CURDIR)/backend/bin
SERVICES_DIR := $(CURDIR)/backend

SERVICES := \
	user

GOLANGCI_LINT_VERSION := v2.12.0
GOLANGCI_CONFIG := $(CURDIR)/backend/.golangci.yaml

ifeq ($(OS),Windows_NT)
	EXE_EXT := .exe
else
	EXE_EXT :=
endif

GOLANGCI_LINT_BIN := $(LOCAL_BIN)/golangci-lint$(EXE_EXT)


.PHONY: run
run:
	docker-compose -f docker-compose.yml up -d --build


.PHONY: stop
stop:
	docker-compose -f docker-compose.yml down -v


.PHONY: generate
generate:
	cd "$(SERVICES_DIR)/user" && go generate ./...


.PHONY: lint
lint: $(GOLANGCI_LINT_BIN)
	@set -eu; \
	echo "Running golangci-lint..."; \
	for service in $(SERVICES); do \
		service_dir="$(SERVICES_DIR)/$$service"; \
		echo " -> $$service"; \
		if [ ! -d "$$service_dir" ]; then \
			echo "Service directory '$$service_dir' not found"; \
			exit 1; \
		fi; \
		( \
			cd "$$service_dir" && \
			"$(GOLANGCI_LINT_BIN)" run \
				--config="$(GOLANGCI_CONFIG)" \
		); \
	done


.PHONY: $(addprefix lint-,$(SERVICES))
$(addprefix lint-,$(SERVICES)): lint-%: $(GOLANGCI_LINT_BIN)
	@if [ ! -d "$(SERVICES_DIR)/$*" ]; then \
		echo "Service directory '$(SERVICES_DIR)/$*' not found"; \
		exit 1; \
	fi
	@echo "Running golangci-lint for $*..."
	@cd "$(SERVICES_DIR)/$*" && \
		"$(GOLANGCI_LINT_BIN)" run \
			--config="$(GOLANGCI_CONFIG)"


$(GOLANGCI_LINT_BIN):
	@echo "golangci-lint not found. Installing..."
	@mkdir -p "$(LOCAL_BIN)"
	@GOBIN="$(LOCAL_BIN)" go install \
		github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)


.PHONY: lint-install
lint-install: $(GOLANGCI_LINT_BIN)


.PHONY: lint-clean
lint-clean:
	@rm -f "$(GOLANGCI_LINT_BIN)"