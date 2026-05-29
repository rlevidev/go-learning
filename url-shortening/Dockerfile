# Build stage
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum (if exists)
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Production stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/main .

RUN chmod +x /app/main

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
