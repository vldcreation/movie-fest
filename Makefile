.PHONY: clean all init generate generate_mocks run all start

.DEFAULT_GOAL := all

all: clean build/main run

build/main: cmd/main.go generate
	@echo "Building..."
	go build -o $@ $<
start:
	@echo "Starting..."
	go run cmd/main.go
run:
	./build/main

clean:
	rm -rf build/*
	rm -rf internal/apis

init: clean generate
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

test_api:
	go clean -testcache
	go test ./tests/...

generate: generated generate_mocks migrateup sqlc

generated: common.api.yaml admin.api.yaml user.api.yaml
	@echo "Generating files..."
	mkdir -p internal/apis/admin internal/apis/common internal/apis/user
	oapi-codegen -config cfg.yaml --package admin -generate types,server,spec admin.api.yaml > internal/apis/admin/api.gen.go
	oapi-codegen -config cfg.yaml --package common -generate types,server,spec common.api.yaml > internal/apis/common/api.gen.go
	oapi-codegen -config cfg.yaml --package user -generate types,server,spec user.api.yaml > internal/apis/user/api.gen.go

INTERFACES_GO_FILES := $(shell find internal/repository -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

.PHONY: migrateup migratedown migrateup1 migratedown1 sqlc
cache?=1
dev?=0
DB_URL?="postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"

ifeq ($(dev), 1)
	DB_URL="postgresql://postgres:secret@localhost:5432/movie_fest?sslmode=disable"
endif

migrateup:
	@echo "Running migrations..."
	@migrate -path db/migrations -database $(DB_URL) -verbose up
migrateup1:
	@echo "Running migrations..."
	@migrate -path db/migrations -database $(DB_URL) -verbose up 1
migratedown:
	@echo "Rolling back migrations..."
	@migrate -path db/migrations -database $(DB_URL) -verbose down
migratedown1:
	@echo "Rolling back migrations..."
	@migrate -path db/migrations -database $(DB_URL) -verbose down 1
sqlc:
	@echo "Generating sqlc..."
	sqlc generate