TEST_PATTERN ?= .
GOBIN := $(shell pwd)/bin
PATH := $(GOBIN):$(PATH)

include .env
export
export GOBIN
export PATH

# Install dependencies
setup:
	go mod download
.PHONY: setup

# Run tests with coverage
test:
	go test $(TEST_OPTIONS) -failfast -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./... -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Display coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# Run linters
lint: bin/golangci-lint
	./bin/golangci-lint run ./...
.PHONY: lint

# Format and organize imports for all Go files
fmt:
	go install golang.org/x/tools/cmd/goimports@v0.5
	find . -name '*.go' -exec gofmt -w -s {} \; -exec goimports -w -local github.com/janschill/track-me {} \;
.PHONY: fmt

# Download and install golangci-lint
bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VERSION)
.PHONY: bin/golangci-lint

# Build the project
build:
	go generate ./...
	go build -o trackme cmd/server/main.go
.PHONY: build

# Build for multiple platforms
build-all-platforms:
	go generate ./...
	env GOOS=darwin  go build -o trackme-darwin  cmd/server/main.go
	env GOOS=linux   go build -o trackme-linux   cmd/server/main.go
	env GOOS=windows go build -o trackme-windows.exe cmd/server/main.go
.PHONY: build-all-platforms

clean:
	@echo "Cleaning up generated files..."
	rm -f trackme trackme-*
	rm -f coverage.*
	rm -rf bin
.PHONY: clean

update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
.PHONY: update-deps

# List TODOs in the codebase
todo:
	@grep --exclude-dir={vendor,node_modules,bin,.git} --exclude=Makefile --text --color -nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

# Start the server in development mode
dev:
	@echo "Starting server"
	go run cmd/server/main.go
.PHONY: dev

dev-watch:
	@echo "Starting server with hot reload..."
	air
.PHONY: dev-watch

# Download the database from a remote server
download-db:
	scp -r $(SSH_USER)@$(SERVER_ADDRESS):$(REMOTE_DB_PATH) .
.PHONY: download-db

upload-db:
	@read -p "Are you sure you want to upload and replace the production DB? This action cannot be undone. Type 'yes' to proceed: " confirm && \
	if [ "$$confirm" = "yes" ]; then \
		scp -r $(DB_PATH) $(SSH_USER)@$(SERVER_ADDRESS):$(REMOTE_DB_PATH)
	else \
		echo "Operation cancelled."; \
	fi
.PHONY: upload-db

# Sync .env file to the remote server
sync-env:
	@read -p "Are you sure you want to sync .env? This action cannot be undone. Type 'yes' to proceed: " confirm && \
	if [ "$$confirm" = "yes" ]; then \
		rsync -avz .env $(SSH_USER)@$(SERVER_ADDRESS):$(REMOTE_APP_ROOT)/.env; \
	else \
		echo "Operation cancelled."; \
	fi
.PHONY: sync-env

# Run the Garmin Outbound mock service
run-mock:
	@echo "Starting Garmin Outbound mock service"
	go run cmd/mockservice/mock_service.go
.PHONY: run-mock

# Database management commands
reset-db create-db destroy-db seed-db clear-db:
	@echo "$(subst -, ,$@)ing database..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="$@"
.PHONY: reset-db create-db destroy-db seed-db clear-db

# Aggregate data for a specific day
aggregate-db:
	@echo "Aggregating data for day $(day)..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="aggregate" -day=$(day)
.PHONY: aggregate-db

# Import GPX file to test folder
import-gpx:
	@echo "Importing GPX file to test folder..."
	@filename=$(file); \
	output="internal/db/test_data/$$(basename $$filename .gpx).json"; \
	go run cmd/gpxparser/main.go $$filename $$output
.PHONY: import-gpx

deploy: clean lint test
	@echo "Building binary using Docker..."
	docker build -t trackme-builder .

	@echo "Creating container to extract binary..."
	docker create --name extract trackme-builder

	@echo "Copying binary from container..."
	docker cp extract:/app/trackme ./trackme

	@echo "Removing container..."
	docker rm extract

	@echo "Copying binary and web directory track-me-temp on server..."
	scp -r trackme web $(SSH_USER)@$(SERVER_ADDRESS):~/track-me-temp/

	@echo "Replacing with new binary and restarting server..."
	ssh $(SSH_USER)@$(SERVER_ADDRESS) 'bash -s' < scripts/restart-server.sh

.PHONY: deploy
