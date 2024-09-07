test: 
	docker compose up -d
	go test ./...

run:
	go run ./cmd/main.go -config ./config.yaml