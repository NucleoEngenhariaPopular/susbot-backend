# Build stage: use an official Golang image with Go 1.23 (or later)
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git (if your modules require it)
RUN apk add --no-cache git

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application source code
COPY . .

# Build the binary.
# The flags "-s -w" strip debugging information for a smaller binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Final stage: use a minimal image
FROM scratch

# Copy certificates for HTTPS (if your APIs need them)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main ./main.go

# Expose the port (adjust if your API listens on a different port)
EXPOSE ${PORT:-8082}

# Run the binary
ENTRYPOINT ["./main.go"]
