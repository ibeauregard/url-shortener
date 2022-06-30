# syntax=docker/dockerfile:1
FROM golang:1.18.3-alpine3.16 as build-image
RUN apk add build-base
COPY ./app/go.mod ./app/go.sum /app/
WORKDIR /app
#RUN go install github.com/mattn/go-sqlite3
RUN go mod download
COPY ./app .
RUN go build -o url_shortener

FROM alpine:3.16.0
WORKDIR /home/urlshortener
COPY --from=build-image /app/url_shortener .
VOLUME ["/home/urlshortener/db"]
WORKDIR /home/urlshortener/db
RUN apk add sqlite
RUN rm -rf /var/cache/apk/*
COPY ./app/db .
WORKDIR /home/urlshortener
EXPOSE 8080
CMD ["/bin/sh", "db/create-db.sh"]
