# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o discordbot ./cmd/discordbot/main.go

# Stage 2: Create a minimal image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/discordbot .

# Expose the port the application listens on (if any)
EXPOSE 5600

# Set environment variables (if any)
# ENV GITHUB_TOKEN=your_github_token

# Command to run the executable
CMD ["./discordbot"]
