.PHONY: up up-db down build logs ps restart \
        run-user-service run-api-gateway \
        install-migrate migrate-up migrate-down migrate-force \
        proto \
        tidy

USER_SERVICE_DB_URL := pgx5://admin:admin@localhost:5432/db?sslmode=disable
USER_SERVICE_MIGRATIONS := services/user-service/internal/migrator/migrations

# ── Docker ───────────────────────────────────────────────────────────────────

up:
	docker compose up -d

up-db:
	docker compose up -d user-service-postgres

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

ps:
	docker compose ps

restart:
	docker compose down && docker compose up -d --build

# ── Local run ────────────────────────────────────────────────────────────────

run-user-service:
	cd services/user-service && go run cmd/main.go --config=config/local.yaml

run-api-gateway:
	cd services/api-gateway && go run cmd/main.go --config=config/local.yaml

# ── Migrations ───────────────────────────────────────────────────────────────

install-migrate:
	go install -tags 'pgx5' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	migrate -path $(USER_SERVICE_MIGRATIONS) -database "$(USER_SERVICE_DB_URL)" up

migrate-down:
	migrate -path $(USER_SERVICE_MIGRATIONS) -database "$(USER_SERVICE_DB_URL)" down 1

migrate-force:
	migrate -path $(USER_SERVICE_MIGRATIONS) -database "$(USER_SERVICE_DB_URL)" force $(version)

# ── Proto ────────────────────────────────────────────────────────────────────

proto-auth:
	protoc \
		--go_out=proto/gen/auth --go_opt=paths=source_relative \
		--go-grpc_out=proto/gen/auth --go-grpc_opt=paths=source_relative \
		--proto_path=proto \
		proto/auth.proto

proto-product:
	protoc \
		--go_out=proto/gen/product --go_opt=paths=source_relative \
		--go-grpc_out=proto/gen/product --go-grpc_opt=paths=source_relative \
		--proto_path=proto \
		proto/product.proto
# ── Go workspace ─────────────────────────────────────────────────────────────

tidy:
	cd pkg && go mod tidy
	cd proto && go mod tidy
	cd services/user-service && go mod tidy
	cd services/api-gateway && go mod tidy