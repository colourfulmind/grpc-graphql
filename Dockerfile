FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go
RUN go build -o client ./cmd/client/main.go
RUN go build -o migrator ./cmd/migrator/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/client /app/client
COPY --from=builder /app/migrator /app/migrator
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/migrations /app/migrations
COPY wait-for-it.sh /app/wait-for-it.sh

EXPOSE 8080
EXPOSE 5432

CMD ["/app/wait-for-it.sh", "db:5432", "--", "/app/migrator", "--config=./configs/config.yml", "/app/server", "--config=./configs/config.yml", "/app/client", "--config=./configs/config.yml"]
