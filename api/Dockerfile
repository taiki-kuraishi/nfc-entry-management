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

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.55.2

# ENV PATH="/root/go/bin:${PATH}"

# lint stage
FROM development AS lint

COPY .golangci.yml .

CMD [ "golangci-lint","run" ]
