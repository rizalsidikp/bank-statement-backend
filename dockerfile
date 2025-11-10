# Use official golang image
FROM golang:1.25 AS builder

WORKDIR /app

# Copy dependency files first for efficient Docker caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o app ./cmd

# Run the built binary
CMD ["./app"]
