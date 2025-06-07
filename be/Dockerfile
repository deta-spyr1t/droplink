# Stage 1: Build
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git (for go modules)
RUN apk add --no-cache git

# Copy go mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -o droplink-be

# Stage 2: Minimal runtime image
FROM scratch

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/droplink-be .

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/droplink-be"]
