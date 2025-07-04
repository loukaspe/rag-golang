FROM golang:1.24.3-alpine3.21

# Install git

RUN apk update && apk add --no-cache git build-base bash

WORKDIR /app

COPY ../go.mod go.sum ./

RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY .. .

RUN go install go.uber.org/mock/mockgen@latest
RUN go install github.com/joho/godotenv/cmd/godotenv@v1.4.0
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
RUN go install github.com/swaggo/swag/cmd/swag@latest