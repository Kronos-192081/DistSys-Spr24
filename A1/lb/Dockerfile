# syntax=docker/dockerfile:1

FROM golang:latest AS build

WORKDIR /lb

RUN apt update && apt install -y docker.io

# Download Go modules
COPY . ./
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-lb

FROM scratch AS deploy

WORKDIR /lb

COPY --from=build /docker-lb /docker-lb

# Expose PORT
EXPOSE 5000

# Run
CMD ["/docker-lb"]

# CMD ["sleep", "infinity"]