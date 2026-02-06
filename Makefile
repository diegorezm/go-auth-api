# ---------- CONFIG ----------

DB_DRIVER=postgres
DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=postgres
DB_PASSWORD?=postgres
DB_NAME?=app

MIGRATIONS_DIR=internal/adapters/postgresql/migrations

DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

GOOSE=goose
SQLC=sqlc

# ---------- HELP ----------

.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "make migrate-up        Run all migrations"
	@echo "make migrate-down      Rollback last migration"
	@echo "make migrate-status    Show migration status"
	@echo "make migrate-reset     Reset database"
	@echo "make migrate-create name=add_users_table"
	@echo ""
	@echo "make sqlc              Generate sqlc code"
	@echo "make db                migrate-up + sqlc"
	@echo ""

# ---------- GOOSE ----------

.PHONY: migrate-up
migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DATABASE_URL)" down

.PHONY: migrate-status
migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DATABASE_URL)" status

.PHONY: migrate-reset
migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DATABASE_URL)" reset

.PHONY: migrate-create
migrate-create:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

# ---------- SQLC ----------

.PHONY: sqlc
sqlc:
	$(SQLC) generate

# ---------- COMBINED ----------

.PHONY: db
db: migrate-up sqlc
