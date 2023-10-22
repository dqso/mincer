FROM golang:1.21.0 AS builder

RUN apt-get update && \
    apt-get install --no-install-recommends --assume-yes \
    protobuf-compiler protoc-gen-go

WORKDIR /app/server

COPY server/go.mod go.mod
COPY server/go.sum go.sum
RUN go mod download

COPY server/pkg pkg
COPY server/internal internal
COPY server/cmd cmd

WORKDIR /app
COPY proto proto

WORKDIR /app/server
RUN go mod tidy
RUN go generate ./...
RUN CGO_ENABLED=0 go build -o ./bin/server ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server/bin/server /app/bin/server

EXPOSE 8080
EXPOSE 12345

ENTRYPOINT ["/app/bin/server"]
