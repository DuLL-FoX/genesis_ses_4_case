# Use the official Golang image as a build environment
FROM golang:1.18-alpine as builder

# Install build tools
RUN apk add --no-cache gcc musl-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Ensure the migrations directory is copied
RUN mkdir -p internal/db/migrations
COPY migrations internal/db/migrations

# Build the Go app
RUN go build -o main ./cmd/app

ENV ETHEREAL_EMAIL=ruthie.beier@ethereal.email
ENV ETHEREAL_PASSWORD=yWRW7aaB4dQMR8Sqsx

# Run tests
RUN go test ./... -v

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
