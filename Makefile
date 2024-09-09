GRAPHIQL_ENABLED = false
DSN = "postgres://user:pass@localhost:5432/testdb?sslmode=disable"

deps:
	docker compose build

run:
	docker compose up -d --no-build

test: 
	docker compose up -d --no-build
	go test ./... -v
	docker compose down

migrate: 
	go run ./migrations/script/migrate.go -path ./migrations/schema.up.sql -dsn $(DSN)
	go run ./migrations/script/migrate.go -path ./migrations/seed.up.sql -dsn $(DSN)
