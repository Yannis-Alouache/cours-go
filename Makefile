APP_NAME := reservation-salles
DIST_DIR := dist
DOCKER_COMPOSE := $(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; elif docker-compose version >/dev/null 2>&1; then echo "docker-compose"; else echo "docker compose"; fi)

.PHONY: run-server run-cli db-up db-down db-logs build build-server-linux build-cli-linux build-cli-darwin build-cli-windows clean release-gh release-snapshot

run-server:
	go run ./cmd/server

run-cli:
	go run ./cmd/cli

db-up:
	$(DOCKER_COMPOSE) up -d postgres

db-down:
	$(DOCKER_COMPOSE) down

db-logs:
	$(DOCKER_COMPOSE) logs -f postgres

build: build-server-linux build-cli-linux build-cli-darwin build-cli-windows

build-server-linux:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/server-linux-amd64 ./cmd/server

build-cli-linux:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/cli-linux-amd64 ./cmd/cli

build-cli-darwin:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/cli-darwin-arm64 ./cmd/cli

build-cli-windows:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/cli-windows-amd64.exe ./cmd/cli

release-gh: build
	gh release create "$$(date +%Y%m%d-%H%M%S)" $(DIST_DIR)/*

release-snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -rf $(DIST_DIR)
