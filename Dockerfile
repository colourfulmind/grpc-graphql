FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/client ./cmd/client/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/client /app/client

COPY --from=builder /app/tests/*.go /app/
COPY --from=builder /app/go.* /app/
COPY --from=builder /app/configs /app/configs

RUN apk add --no-cache postgresql-client go --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing grpcurl

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./server", "--config=./configs/config.yml"]
