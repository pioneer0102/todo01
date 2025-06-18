FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Install required Go tools
RUN go install github.com/volatiletech/sqlboiler/v4@latest && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest && \
    go install github.com/bufbuild/buf/cmd/buf@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN mkdir -p bin
RUN go build -o bin/server ./cmd/server

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/server /app/server

# Run the application
CMD ["./server"] 