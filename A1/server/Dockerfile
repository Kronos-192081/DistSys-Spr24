# syntax=docker/dockerfile:1

FROM golang:latest AS build
# LABEL maintainer "Kartik Pontula"

WORKDIR /server

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-server

FROM scratch AS deploy

WORKDIR /server

COPY --from=build /docker-server /docker-server

# Expose PORT
EXPOSE 5000

# Run
ENTRYPOINT ["/docker-server"]