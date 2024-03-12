# syntax=docker/dockerfile:1

FROM golang:latest AS build
# RUN apt update && apt install -y docker.io
WORKDIR /lb
# COPY lb.go ./
# COPY conhash ./conhash
# RUN go mod init distri-lb && go mod edit -replace prakhar/conhash=./conhash && go mod tidy && go mod download
COPY . ./
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o distri-lb

FROM ubuntu AS deploy
RUN apt-get -y update && apt-get -y upgrade
RUN apt-get install -y sqlite3 libsqlite3-dev
RUN /usr/bin/sqlite3 lb.db
# CMD /bin/bash
COPY --from=build /lb/distri-lb ./
EXPOSE 5000
CMD ["/distri-lb"]
# CMD ["sleep", "infinity"]