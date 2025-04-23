# Prompt Copilot: Dockerfile for a Go application
# This Dockerfile is designed to build and run a Go application  
# Start from a Go base image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the source code from backend folder
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o fluffy-enigma .

# Use a small alpine image for the final container
FROM alpine:3.21

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/fluffy-enigma .

# Run the application
CMD ["./fluffy-enigma"]