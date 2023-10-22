.PHONY: run

run: generate
	POSTGRES_PASSWORD_FILE=config/pg_password \
		NC_PRIVATE_KEY_FILE=config/nc_private_key \
#		LOG_LEVEL=debug \
		LOG_LEVEL=info \
		go run ./cmd/main.go

generate:
	go generate ./...
