SHELL := /bin/bash

all: users-api metrics

keys:
	go run ./cmd/users-admin/main.go keygen private.pem

admin:
	go run ./cmd/users-admin/main.go --db-disable-tls=1 useradd admin@example.com gophers

migrate:
	go run ./cmd/users-admin/main.go --db-disable-tls=1 migrate

seed: migrate
	go run ./cmd/users-admin/main.go --db-disable-tls=1 seed

users-api:
	docker build \
		-f dockerfile.users-api \
		-t igorgomonov/users-api-amd64:1.0 \
		--build-arg PACKAGE_NAME=users-api \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

metrics:
	docker build \
		-f dockerfile.metrics \
		-t igorgomonov/users-metrics-amd64:1.0 \
		--build-arg PACKAGE_NAME=metrics \
		--build-arg PACKAGE_PREFIX=sidecar/ \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

up:
	docker-compose up

down:
	docker-compose down

test:
	go test -mod=vendor ./... -count=1

clean:
	docker system prune -f

stop-all:
	docker stop $(docker ps -aq)

remove-all:
	docker rm $(docker ps -aq)

tidy:
	go mod tidy
	go mod vendor

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -t -d -v ./...

deps-cleancache:
	go clean -modcache