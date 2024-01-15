# syntax=docker/dockerfile:1

FROM golang:latest

# LABEL maintainer "Kartik Pontula"

# Expose PORT
EXPOSE 5000

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-server

# Run
CMD ["/docker-server"]