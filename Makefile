include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/pkg/errors
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/joho/godotenv
	go get -u github.com/jackc/pgconn
	go get -u go.uber.org/zap
	go get -u github.com/fatih/color
	go get -u github.com/grpc-ecosystem/go-grpc-middleware
	go get -u github.com/natefinch/lumberjack
	go get -u github.com/georgysavva/scany/pgxscan
	go get -u github.com/Masterminds/squirrel
	go get -u gopkg.in/gomail.v2\
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2

generate:
	make generate-access-api
	make generate-auth-api

generate-access-api:
	mkdir -p pkg/access
	protoc --proto_path api/access \
		--go_out=pkg/access --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/access --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		api/access/access.proto

generate-auth-api:
	mkdir -p pkg/auth
	protoc --proto_path=api/auth --proto_path=vendor.protogen \
		--go_out=pkg/auth --go_opt=paths=source_relative \
		--go-grpc_out=pkg/auth --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pkg/auth --grpc-gateway_opt=paths=source_relative,logtostderr=true \
		--validate_out="lang=go:pkg/auth" --validate_opt=paths=source_relative \
		api/auth/auth.proto



vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
        			mkdir -p vendor.protogen/validate &&\
        			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
        			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
        			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi


create-database:
	PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -U $(PG_USER) -p $(PG_PORT) -c "CREATE DATABASE $(PG_DATABASE_NAME);"

drop-database:
	PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -U $(PG_USER) -p $(PG_PORT) -c "DROP DATABASE IF EXISTS $(PG_DATABASE_NAME);"

local-migrations-status:
	${LOCAL_BIN}/goose -dir $(LOCAL_MIGRATION_DIR) postgres ${PG_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} down -v

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
