# ==========================================
# Stage 1: Build the Go binary
# ==========================================
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install build dependencies if needed (e.g., git, ca-certificates)
RUN apk add --no-cache git ca-certificates

# Copy dependency manifests first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies; this layer is cached unless go.mod or go.sum changes
RUN go mod download && go mod verify

# Copy the rest of the application source code
COPY . .

# Build a statically linked, stripped production binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/server .

# ==========================================
# Stage 2: Create a minimal runtime image
# ==========================================
FROM alpine:3.19 AS runner

# Create a non-root user and group for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy CA certificates for secure network outgoing calls (HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server /app/server

# Use the non-root user to run the application
USER appuser

# Expose your application's network port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/server"]
