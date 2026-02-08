# === Stage 1: Build Frontend ===
# Use --platform=$BUILDPLATFORM to always run the build on the efficient native architecture (e.g. amd64)
# The output (HTML/JS/CSS) is architecture-agnostic.
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

# Install dependencies first for better caching
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

# Copy sources and build
COPY frontend/ .
RUN npm run build

# === Stage 2: Build Backend ===
# Run the Go compiler on the native architecture for speed (Cross-Compilation)
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend

# Install build dependencies
RUN apk add --no-cache git

# Copy Go modules manifests
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy frontend build artifacts to backend embedding directory
COPY --from=frontend-builder /app/frontend/build ./internal/ui/dist

# Copy backend source code
COPY backend/ .

# Build the binary using cross-compilation
# TARGETOS and TARGETARCH are automatically set by docker buildx
ARG TARGETOS
ARG TARGETARCH

# -s -w: Strip debug information
# CGO_ENABLED=0: Static binary
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -trimpath -o /app/wolite main.go

# Create data directory with correct permissions
RUN mkdir /app/data

# === Stage 3: Final Minimal Image ===
# Use distroless for security and minimal size
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /

# Copy the binary from the backend builder
COPY --from=backend-builder /app/wolite /wolite

# Copy /data directory with correct ownership
COPY --from=backend-builder --chown=65532:65532 /app/data /data

# Set environment variables
ENV DATABASE_PATH=/data/wolite.json
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run as non-root user (provided by distroless)
USER 65532:65532

# Set the entrypoint
ENTRYPOINT ["/wolite"]
