# Use the official Golang image to create a build artifact.
FROM golang:1.15 as builder

# Copy local code to the container image.
WORKDIR /app
COPY . ./

# Build the command inside the container.
# (Ensure that the module is executable and not dependent on any code generated during the build process.)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o hanoi

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/hanoi /hanoi

# Run the web service on container startup.
CMD ["/hanoi"]
