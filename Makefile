.PHONY: reset-db create-db destroy-db

DB_PATH="./data/trips.db"

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
