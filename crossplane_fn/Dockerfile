# Use an Alpine-based Golang image for both building and running the application
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the entire source directory to the container
COPY . .

# Build the application to a binary named "hanoi-function"
# This command builds the app assuming it has no external dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -o hanoi-function .

# Start a new stage from scratch to keep the image small
FROM alpine:latest  

# Install ca-certificates in case the application makes outgoing HTTPS calls
RUN apk --no-cache add ca-certificates

# Set the working directory in the container to /root/
WORKDIR /root/

# Copy the binary from the builder stage to the current stage
COPY --from=builder /app/hanoi-function .

# Expose port 8080 on the container
EXPOSE 8080

# Command to run the executable
CMD ["./hanoi-function"]
