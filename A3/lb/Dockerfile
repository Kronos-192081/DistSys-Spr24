# syntax=docker/dockerfile:1

FROM golang:latest AS build
# RUN apt update && apt install -y docker.io
# RUN apt install -y curl
WORKDIR /lb
COPY lb.go ./
COPY conhash ./conhash
RUN go mod init distri-lb && go mod edit -replace prakhar/conhash=./conhash && go mod tidy && go mod download
# COPY . ./
# RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o distri-lb

# FROM ubuntu AS deploy
# RUN apt-get -y update && apt-get -y upgrade
# RUN apt-get install -y sqlite3 libsqlite3-dev
# RUN /usr/bin/sqlite3 lb.db

# COPY --from=build /lb/distri-lb ./
# EXPOSE 5000
# CMD ["/distri-lb"]

FROM postgres:latest AS deploy
COPY --from=build /lb/distri-lb ./
COPY init.sh /docker-entrypoint-initdb.d/

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 20CS30061
ENV POSTGRES_DB testdb

EXPOSE 5000