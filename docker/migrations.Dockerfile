FROM golang:1.21.0-alpine

RUN go install github.com/jackc/tern/v2@latest

COPY ./server/migrations /migrations

CMD tern version && \
    echo "\$DESTINATION=$DESTINATION" && \
    echo && \
    echo 'check a status before migration...' && \
    tern status \
        --host postgres \
        --config /run/secrets/migrations_config \
        --migrations /migrations && \
    echo && \
    tern migrate \
        --host postgres \
        --config /run/secrets/migrations_config \
        --migrations /migrations \
        --destination $DESTINATION || \
    # using a double pipe || to check the status after a failed migration.
    echo && \
    echo 'check a status after migration...' && \
    tern status \
        --host postgres \
        --config /run/secrets/migrations_config \
        --migrations /migrations
