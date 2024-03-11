# syntax=docker/dockerfile:1

FROM golang:latest AS build
RUN apt update && apt install -y docker.io
WORKDIR /lb
COPY lb.go ./
RUN go mod init distri-lb && go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o distri-lb

FROM scratch AS deploy
COPY --from=build /lb/distri-lb ./
EXPOSE 5000
ENTRYPOINT ["/distri-lb"]