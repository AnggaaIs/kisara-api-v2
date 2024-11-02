# Build stage
FROM golang:1.23.2-alpine3.20 AS builder

LABEL name "kisara-api (build stage)"
LABEL maintainer "AnggaaIs <servantangga@gmail.com>"

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src/main.go

# Final stage
FROM alpine:latest

LABEL name "kisara-api (final stage)"
LABEL maintainer "AnggaaIs <servantangga@gmail.com>"

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Set the executable permission
RUN chmod +x /app/main

# Expose port
EXPOSE 3000

# Command to run the application
CMD ["./main"]