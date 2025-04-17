# ========================
# Ingest-gateway Dev Makefile
# ========================

# Note: Assumes a Debian-based environment for tooling installs

# App entrypoint
MAIN_PKG := ./main.go

# Build, tag, and version vars
VERSION_FILE := VERSION
VERSION ?= $(shell cat $(VERSION_FILE))
IMAGE_REPO ?= 192.168.2.17:5000
IMAGE_NAME ?= embedding-service
IMAGE_TAG ?= $(VERSION)
IMAGE_FULL := $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

# =================================
# Dev Quality of Life
# =================================

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: run
run:
	go run $(MAIN_PKG)

.PHONY: lint
lint:
	golangci-lint run

# =================================
# Docker build, tag, and push
# =================================

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_FULL) .

.PHONY: docker-push
docker-push:
	docker push $(IMAGE_FULL)

.PHONY: tag
tag:
	@version=$$(cat VERSION); \
	git tag $$version; \
	git push origin main --tags