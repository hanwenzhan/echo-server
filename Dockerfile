# syntax=docker/dockerfile:1.4

# Build App Binary
FROM golang:alpine AS build
ADD . /src/
WORKDIR /src
RUN go mod tidy
RUN go build -o /build/main main.go

# Build Docker Image
FROM alpine 
LABEL Description="line echo server"

#RUN apk add --no-cache ca-certificates

COPY .env /app/.env
COPY --from=build /build/main /app/main
WORKDIR /app
ENTRYPOINT ["./main"]
