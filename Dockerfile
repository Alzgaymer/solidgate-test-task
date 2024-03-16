# Start from the official Golang image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/server

# Start a new stage from scratch
FROM alpine:latest

# Set a non-root user
RUN adduser -D -g '' appuser

# Set the Current Working Directory inside the container
WORKDIR /home/appuser

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/app .

# Change the owner of the binary to the non-root user
RUN chown appuser:appuser ./app

# Expose port 8080 to the outside world
EXPOSE 8080

# Switch to the non-root user
USER appuser

# Command to run the executable
CMD ["./app"]
