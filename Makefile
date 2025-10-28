# --------
# Makefile
# --------
SHELL := /bin/bash

APP        ?= pdfgo
MAIN       ?= ./cmd/$(APP)     # chemin du package main (change si autre)
PKG        ?= ./...
BIN_DIR    ?= bin
VENV_DIR   ?= .venv
PY         ?= python3

GO         ?= go
GOFLAGS    ?=
TESTFLAGS  ?= -race -covermode=atomic -coverprofile=coverage.out

# Résout GOBIN même s'il n'est pas défini
GOBIN := $(shell $(GO) env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell $(GO) env GOPATH)/bin
endif

PRE_COMMIT ?= $(VENV_DIR)/bin/pre-commit
STATICCHECK := $(GOBIN)/staticcheck

# ----------------
# Cibles principales
# ----------------
.PHONY: help setup hooks fmt fmt-check vet lint test test-short cover cover-html tidy build clean pre-commit-run

help:
	@echo "Cibles disponibles :"
	@echo "  setup           - Prépare l'env: venv + pre-commit + staticcheck + go mod download"
	@echo "  hooks           - Installe les hooks pre-commit"
	@echo "  fmt             - gofmt -s -w ."
	@echo "  fmt-check       - Vérifie que le code est formaté"
	@echo "  vet             - go vet ./..."
	@echo "  lint            - staticcheck ./..."
	@echo "  test            - go test (race+coverage)"
	@echo "  test-short      - go test -short"
	@echo "  cover           - Résumé couverture (coverage.out)"
	@echo "  cover-html      - Rapport HTML (coverage.html)"
	@echo "  tidy            - go mod tidy"
	@echo "  build           - Build l'app dans $(BIN_DIR)/$(APP)"
	@echo "  clean           - Nettoie bin/ et artefacts (coverage.*)"
	@echo "  pre-commit-run  - Exécute tous les hooks sur tous les fichiers"

# Prépare l'environnement out-of-the-box (sans toucher au Python système)
setup: $(VENV_DIR)/bin/activate $(PRE_COMMIT) $(STATICCHECK)
	@echo "==> go mod download"
	$(GO) mod download
	@echo "==> installation des hooks pre-commit"
	$(PRE_COMMIT) install
	@echo "==> setup OK"

# Active/installe le venv local
$(VENV_DIR)/bin/activate:
	@echo "==> création du venv $(VENV_DIR)"
	$(PY) -m venv $(VENV_DIR)
	@echo "==> mise à jour de pip"
	$(VENV_DIR)/bin/pip install -U pip

# Installe pre-commit dans le venv
$(PRE_COMMIT): $(VENV_DIR)/bin/activate
	@echo "==> installation de pre-commit dans $(VENV_DIR)"
	$(VENV_DIR)/bin/pip install pre-commit

# Installe staticcheck côté Go
$(STATICCHECK):
	@echo "==> installation de staticcheck"
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

hooks:
	$(PRE_COMMIT) install

fmt:
	@gofmt -s -w .

fmt-check:
	@echo "==> vérification du formatage (gofmt)"
	@test -z "$$(gofmt -s -l . | grep -v '^vendor/' | tee /dev/stderr)" || (echo "✗ Code non formaté. Lance 'make fmt'."; exit 1)

vet:
	$(GO) vet $(PKG)

lint: $(STATICCHECK)
	$(STATICCHECK) $(PKG)

test:
	$(GO) test $(TESTFLAGS) $(PKG)

test-short:
	$(GO) test -short $(PKG)

cover:
	@$(GO) tool cover -func=coverage.out

cover-html:
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Rapport: coverage.html"

tidy:
	$(GO) mod tidy

# Build binaire (ne commitera rien si ton .gitignore ignore bin/)
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP) $(MAIN)
	@echo "Binaire: $(BIN_DIR)/$(APP)"

clean:
	@rm -rf $(BIN_DIR) coverage.out coverage.html

pre-commit-run: $(PRE_COMMIT)
	$(PRE_COMMIT) run --all-files
