deps:
	docker compose up -d

test: deps 
	go test ./...

migrate: deps
	go run ./migrations/script/migrate.go -path ./migrations/schema.up.sql -dsn "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
	go run ./migrations/script/migrate.go -path ./migrations/seed.up.sql -dsn "postgres://user:pass@localhost:5432/testdb?sslmode=disable"

run: deps
	go run ./cmd/main.go -config ./config.yaml -graphiql