FROM golang:1.24.3-alpine3.21 as builder

WORKDIR /app

COPY ../go.mod go.sum ./
RUN go mod download
RUN go install github.com/joho/godotenv/cmd/godotenv@v1.4.0
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM golang:1.24.3-alpine3.21

RUN apk update
RUN apk add build-base bash

COPY --from=builder /go /go

WORKDIR /app

COPY .. .

RUN GOOS=linux go build -gcflags='all=-N -l' -tags musl -a -installsuffix cgo -o main ./cmd/http/main.go

EXPOSE 8080
EXPOSE 40000

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "--continue=true", "exec", "main"]