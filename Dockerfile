# Use the official Golang image as the base image
FROM golang:1.25.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the bot binary
RUN go build -o adhd-bot cmd/bot/main.go

# Build the API binary
RUN go build -o adhd-api cmd/api/main.go

# Build the migration tool
RUN go build -o migrate cmd/migrate/main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binaries from the builder stage
COPY --from=builder /app/adhd-bot .
COPY --from=builder /app/adhd-api .
COPY --from=builder /app/migrate .

# Copy migration files
COPY --from=builder /app/internal/infra/postgres/migrations ./internal/infra/postgres/migrations

# Expose the port the API runs on
EXPOSE 8080

# Command to run the executable
CMD ["./adhd-bot"]