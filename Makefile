.PHONY: migrate
migrate:
	docker-compose run --rm app ./migrator

.PHONY: build_server
build_server:
	docker-compose run --rm app go run ./server --config=./configs/config.yml

.PHONY: build_client
build_client:
	docker-compose run --rm app go run ./client --config=./configs/config.yml

.PHONY: tests
tests:
	docker-compose run --rm app cd tests && go test -v
