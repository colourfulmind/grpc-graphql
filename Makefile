.PHONY: build_server
build_server:
	go build -o server cmd/server/main.go
	./server --config=./configs/config.yml

.PHONY: build_client
build_client:
	go build -o client cmd/client/main.go
	./client --config=./configs/config.yml

.PHONY: tests
tests:
	cd tests && go test -v

.PHONY: clean
clean:
	rm -rf server client

generate_protos:
	protoc -I protos/proto protos/proto/*/*.proto --go_out=./protos/gen/go --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go --go-grpc_opt=paths=source_relative

generate_graphql:
	cd internal/client && go get github.com/99designs/gqlgen@v0.17.47 && go run github.com/99designs/gqlgen generate
