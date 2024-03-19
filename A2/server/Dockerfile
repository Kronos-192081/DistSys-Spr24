#syntax=docker/dockerfile:1

# build stage
FROM golang:latest AS build
WORKDIR /server
COPY main.go ./
RUN go mod init distri-server && go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o distri-server

# deploy stage
FROM postgres:latest AS deploy
# WORKDIR /docker-entrypoint-initdb.d/
COPY --from=build /server/distri-server ./
COPY init.sh /docker-entrypoint-initdb.d/
# RUN chmod +x /docker/entrypoint-initdb.d/init.sh

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 20CS30061
ENV POSTGRES_DB testdb

# Expose PORT
EXPOSE 5000

# Run
# CMD ["./init.sh"]
