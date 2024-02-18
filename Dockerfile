# Build stage
FROM golang:alpine3.19 AS build

WORKDIR /app

COPY . .

RUN go build && go test -v ./internal/
# Minimal image with Alpine Linux
FROM alpine:latest

WORKDIR /app

# Copy files from the build stage
COPY --from=build /app/dockerrestapi .

EXPOSE 8081

# Run the program
CMD ["./dockerrestapi"]