# syntax=docker/dockerfile:1
FROM golang:1.18.3-alpine3.16 as build-image
RUN apk add --no-cache build-base
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download
RUN go install github.com/mattn/go-sqlite3
COPY . .
RUN APP_HOST=apphost go test -v -coverprofile=coverage.txt ./...
RUN go build -o url_shortener ./internal

FROM alpine:3.16.0
WORKDIR /home/urlshortener
COPY --from=build-image /app/url_shortener /app/coverage.txt ./
VOLUME ["/home/urlshortener/db/data"]
VOLUME ["/home/urlshortener/tests/unit"]
RUN apk add --no-cache sqlite
COPY internal/db ./db
COPY internal/templates ./templates
COPY internal/static ./static
EXPOSE 8080
CMD ["/bin/sh", "db/setup.sh"]
