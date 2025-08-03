include .envrc

MIGRATIONS_PATH=cmd/migrate/migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))



.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@echo "Running migrations down..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" down $(filter-out $@,$(MAKECMDGOALS))