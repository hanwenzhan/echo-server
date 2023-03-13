# syntax=docker/dockerfile:1.4

# Build App Binary
FROM golang:alpine AS build
ADD . /src/
WORKDIR /src
RUN go mod tidy
RUN go build -o /build/main main.go

# Build Docker Image
FROM scratch 
LABEL Description="go echo server"

COPY --link --from=build /build/main /app/main
ENTRYPOINT ["/app/main"]
