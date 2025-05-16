.PHONY: build
build:
	go build -o build/bin cmd/app/main.go

.PHONY: run
run: build
	build/bin

.PHONY: mock
mock:
	go run go.uber.org/mock/mockgen@latest \
		-source internal/service/service.go \
		-destination internal/service/mocks/mocks.go \
		-package mocks


.PHONY: docs
docs:
	swag init -g ./cmd/app/main.go -o ./docs --parseDependency --parseInternal

.PHONY: gen
gen: mock docs

.PHONY: test
test: unit-test integration-test

.PHONY: dev
dev:
	docker compose up --build postgres service

.PHONY: full-compose-up
full-compose-up:
	docker-compose up --build -d 

.PHONY: migrate
migrate: 
	migrate -source file://schema/migrations -database postgres://postgres:5432@127.0.0.1:5432/tachki_crm\?sslmode=disable up

.PHONY: compose-migrate
compose-migrate:
	docker compose up --build migrations

.PHONY: proto
proto: 
	protoc --go_out=api --go_opt=paths=source_relative \
       --go-grpc_out=api --go-grpc_opt=paths=source_relative \
       proto/proto.proto
