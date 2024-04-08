include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres redis && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### remove docker volume
	docker volume rm shortlinkapi_pg-data
.PHONY: docker-rm-volume

docker-it-db:
	docker exec -it postgres psql -U CodeMaster482 -d FLibraryDB 
.PHONY: docker-it-db

swag-v1: ### swag init
	swag init -g ./internal/app/app.go
.PHONY: swag-v1

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

lint: linter-golangci linter-hadolint linter-dotenv ### run all linters
.PHONY: lint

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

mock: ### run mockgen
	~/go/bin/mockgen -source=./internal/usecase/link.go -destination=./internal/usecase/mocks/mocks.go
	~/go/bin/mockgen -source=./internal/repository/postgres/postgres.go -destination=./internal/repository/postgres/mocks/mocks.go
	~/go/bin/mockgen -source=./internal/delivery/http/handler/handler.go -destination=./internal/delivery/http/handler/mocks/mocks.go
.PHONY: mock

easyjson: ### run easyjson generation
	~/go/bin/easyjson -all internal/model/link.go
	~/go/bin/easyjson -all internal/delivery/http/dto/link.go
.PHONY: easyjson

protoc: ### run protoc generation
	protoc --go_out=../internal/delivery/grpc/generated/ --go-grpc_out=../internal/delivery/grpc/generated/ --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative link.proto
.PHONY: protoc

bin-dep:
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(LOCAL_BIN) go install github.com/mailru/easyjson/...@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
.PHONY: bin-dep

coverage:
	sh scripts/coverage.sh
.PHONY: coverage
