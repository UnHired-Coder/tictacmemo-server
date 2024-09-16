HIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: help build init shell runserver
help:
	make -pRrq  -f $(THIS_FILE) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

build:
	docker-compose build

runserver:
	docker-compose up

shell:
	docker exec -it game-server sh

logs:
	docker logs game-server

test:
	docker exec -it game-server sh -c 'go clean -testcache && ENV=test go test -race -p 1 -v ./...'

down:
	docker-compose down

test-go-acc:
	docker exec -it game-server sh -c 'go clean -testcache && ENV=test go-acc --covermode atomic --output coverage.out ./... -- -race -v -p 1'
