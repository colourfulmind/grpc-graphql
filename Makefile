.PHONY: migrate
migrate:
	go run cmd/migrator/main.go

.PHONY: build_server
build_server:
	go run cmd/server/main.go --config=./configs/config.yml

.PHONY: build_client
build_client:
	go run cmd/client/main.go --config=./configs/config.yml

.PHONY: tests
tests:
	cd tests && go test -v

.PHONY: generate_protos
generate_protos:
	protoc -I protos/proto protos/proto/*/*.proto --go_out=./protos/gen/go --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go --go-grpc_opt=paths=source_relative

.PHONY: generate_graphql
generate_graphql:
	cd internal/client && go get github.com/99designs/gqlgen@v0.17.47 && go run github.com/99designs/gqlgen generate
