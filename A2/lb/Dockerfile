# syntax=docker/dockerfile:1

FROM golang:latest AS build
RUN apt update && apt install -y docker.io
WORKDIR /lb
# COPY . ./
COPY lb.go ./
COPY conhash ./conhash
RUN go mod init distri-lb && go mod edit -replace prakhar/conhash=./conhash && go mod tidy && go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o distri-lb

FROM ubuntu AS deploy
COPY --from=build /lb/distri-lb ./
EXPOSE 5000
CMD ["/distri-lb"]
# CMD ["sleep", "infinity"]