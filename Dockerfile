# ==========================================
# Stage 1: Build the Go binary (Native Speed)
# ==========================================
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

# Create the non-root user details here using your native CPU
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Target architecture automatically populated by buildx
ARG TARGETARCH

# Cross-compile the Go binary to ARM64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/hello hello.go

# ==========================================
# Stage 2: Create a minimal runtime image
# ==========================================
# scratch has no OS files, no shell, and needs zero emulation
FROM scratch AS runner

# Copy the passwd and group files we made in Stage 1
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy CA certificates for secure network calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled ARM64 binary
COPY --from=builder /app/hello /app/hello

# Use the non-root user
USER appuser

EXPOSE 8080 502

ENTRYPOINT ["/app/hello"]
