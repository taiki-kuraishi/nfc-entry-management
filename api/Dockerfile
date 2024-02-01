FROM golang:1.21.6-alpine

WORKDIR /app

COPY ./src .

RUN apk add --no-cache git && \
    go install -v golang.org/x/tools/cmd/goimports@latest && \
    go install -v golang.org/x/tools/gopls@latest

RUN go mod tidy
