#syntax=docker/dockerfile:1

# build stage
FROM golang:latest AS build
WORKDIR /server
COPY main.go ./
RUN go mod init distri-server && go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /distri-server

# deploy stage
FROM scratch AS deploy
WORKDIR /server
COPY --from=build /distri-server /distri-server

# Expose PORT
EXPOSE 5000

# Run
ENTRYPOINT ["/distri-server"]