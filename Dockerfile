# Build the app
FROM golang:1.21.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

# Minimize the image
FROM alpine:latest

COPY --from=builder /app/main /app/main
COPY config.yaml /app/config.yaml

ENV GRAPHIQL_ENABLED=false

EXPOSE 8080

CMD ["sh", "-c", "/app/main -config /app/config.yaml -graphiql=${GRAPHIQL_ENABLED}"]