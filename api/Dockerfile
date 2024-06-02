# syntax=docker/dockerfile:1

# Build the binary
# Start from the latest golang base image
FROM golang:1.22.3 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . .

# Build the binary
RUN go build -o bin .

# The final image

# Using debian as a smaller base image
FROM debian:stable-slim AS deployed

WORKDIR /

COPY --from=builder /app/bin /app/bin

EXPOSE 8080

ENTRYPOINT [ "/app/bin" ]
