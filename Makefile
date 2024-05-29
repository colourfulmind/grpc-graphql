.PHONY: build_server
build_server:
	go build cmd/server/main.go
	./main --config=./configs/config.yml

.PHONY: build_client
build_client:
	go build cmd/client/main.go
	./main --config=./configs/config.yml

.PHONY: clean
clean:
	rm -rf main

generate_files:
	protoc -I protos/proto protos/proto/*/*.proto --go_out=./protos/gen/go --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go --go-grpc_opt=paths=source_relative
