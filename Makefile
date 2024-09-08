GRAPHIQL_ENABLED = false
DSN = "postgres://user:pass@localhost:5432/testdb?sslmode=disable"

deps:
	docker compose up -d

test: deps 
	go test ./...

migrate: deps
	go run ./migrations/script/migrate.go -path ./migrations/schema.up.sql -dsn $(DSN)
	go run ./migrations/script/migrate.go -path ./migrations/seed.up.sql -dsn $(DSN)

build:
	docker build -t app .
	
run: deps
	docker run -p 8080:8080 -e GRAPHIQL_ENABLED=$(GRAPHIQL_ENABLED) app