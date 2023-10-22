.PHONY: build run_server run_client generate_server generate_client

run_server: generate_server
	cd server; POSTGRES_PASSWORD_FILE=config/pg_password \
		NC_PRIVATE_KEY_FILE=config/nc_private_key \
		LOG_LEVEL=info \
		go run ./cmd/main.go

generate_server:
	cd client; go generate ./...

run_client: generate_client
	cd client; go run ./cmd/main.go

generate_client:
	cd client; go generate ./...

build: generate_server generate_client
	cd client; go build -o ../bin/client ./cmd/main.go
	cd server; go build -o ../bin/server ./cmd/main.go

#build_docker: generate_server generate_client
#	docker build -f ./docker/client_build.Dockerfile --output bin .