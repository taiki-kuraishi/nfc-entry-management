FROM golang:1.21.6-alpine AS base

WORKDIR /app

COPY ./src .


RUN apk add --no-cache git && \
    go install -v golang.org/x/tools/cmd/goimports@latest && \
    go install -v golang.org/x/tools/gopls@latest

RUN go mod tidy

ENV PATH="/app/bin:${PATH}"

# development stage
FROM base AS development

WORKDIR /app

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.55.2

# lint stage
FROM development AS lint

WORKDIR /app

CMD [ "golangci-lint","run","--config",".golangci.yml" ]
