DSN = postgres://user:pass@localhost:5432/testdb?sslmode=disable

run:
	docker compose up -d

stop:
	docker compose down

test: 
	@make run 
	go test ./... -v
	@make stop

migrate: 
	go run ./migrations/script/migrate.go -path ./migrations/schema.up.sql -dsn $(DSN)
	go run ./migrations/script/migrate.go -path ./migrations/seed.up.sql -dsn $(DSN)
