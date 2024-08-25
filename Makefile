SOURCE_FILES?=./...
TEST_PATTERN?=.
include .env
export

export GOBIN := $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

setup:
	go mod download
.PHONY: setup

test:
	go test $(TEST_OPTIONS) -failfast -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# Run all the linters
lint: bin/golangci-lint
	./bin/golangci-lint run ./...
.PHONY: lint

# gofmt and goimports all go files
fmt:
	go install golang.org/x/tools/cmd/goimports@v0.5
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w -local github.com/stripe/stripe-cli "$$file"; done
.PHONY: fmt

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VERSION)

build:
	go generate ./...
	go build -o trackme cmd/server/main.go
.PHONY: build

build-all-platforms:
	go generate ./...
	env GOOS=darwin go build -o trackme-darwin cmd/server/main.go
	env GOOS=linux go build -o trackme-linux cmd/server/main.go
	env GOOS=windows go build -o trackme-windows.exe cmd/server/main.go
.PHONY: build-all-platforms

todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=bin \
		--exclude-dir=.git \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

dev:
	@echo "Starting server"
	go run cmd/server/main.go

download-db:
	scp -r $(SSH_USER)@$(SERVER_ADDRESS):$(REMOTE_DB_PATH) .
.PHONY: download-db

sync-env:
	@read -p "Are you sure you want to sync .env? This action cannot be undone. Type 'yes' to proceed: " confirm && \
	if [ "$$confirm" = "yes" ]; then \
		rsync -avz .env $(SSH_USER)@$(SERVER_ADDRESS):$(REMOTE_APP_ROOT)/.env; \
	else \
		echo "Operation cancelled."; \
	fi
.PHONY: sync-env

run-mock:
	@echo "Starting Garmin Outbound mock service"
	go run cmd/mockservice/mock_service.go

reset-db:
	@echo "Resetting database..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="reset"

create-db:
	@echo "Creating database..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="create"

destroy-db:
	@echo "Destroying database..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="destroy"

seed-db:
	@echo "Seeding database..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="seed"

clear-db:
	@echo "Clearing database..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="clear"

aggregate-db:
	@echo "Aggregating data for day $(day)..."
	go run cmd/db/main.go -dbpath=$(DB_PATH) -operation="aggregate" -day=$(day)

import-gpx:
	@echo "Importing GPX file to test folder..."
	@filename=$(file); \
	output="internal/db/test_data/$$(basename $$filename .gpx).json"; \
	go run cmd/gpxparser/main.go $$filename $$output
