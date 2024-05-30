.PHONY: build_client
build:
	docker exec -it ozon-app-1 ./client --config=./configs/config.yml

.PHONY: tests
tests:
	docker exec -it ozon-app-1 go test -v
