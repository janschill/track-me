.PHONY: reset-db create-db destroy-db

DB_PATH="./data/trips.db"

run:
	@echo "Starting server"
	go run cmd/server/main.go

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
