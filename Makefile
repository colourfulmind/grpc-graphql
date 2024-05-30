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
	docker-compose run --rm app cd tests && go test -v
