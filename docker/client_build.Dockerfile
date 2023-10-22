FROM golang:1.21.0 AS build-stage

RUN apt-get update && \
    apt-get install --no-install-recommends --assume-yes \
    protobuf-compiler protoc-gen-go \
    gcc \
    libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev \
    libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

WORKDIR /app/client

COPY client/go.mod go.mod
COPY client/go.sum go.sum
RUN go mod download

COPY client/internal internal
COPY client/fonts fonts
COPY client/cmd cmd

WORKDIR /app
COPY proto proto

WORKDIR /app/client
RUN go mod tidy
RUN go generate ./...
RUN go build -o /app/bin/client ./cmd/main.go

FROM scratch AS export-stage

COPY --from=build-stage /app/bin/client /client
