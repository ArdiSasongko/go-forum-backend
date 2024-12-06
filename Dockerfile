# Build stage
FROM golang:1.23.2-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy application files
COPY . ./

# Build the application
RUN go build -o cmd/bin/main cmd/main.go

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/cmd/bin/main ./cmd/bin/main

# Copy necessary files
COPY --from=builder /app/.env .env
COPY --from=builder /app/root.crt ./root.crt

# Expose application port
EXPOSE 8080

# Set the default command to run the application
CMD ["./cmd/bin/main"]
