# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git curl libstdc++

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install tailwindcss standalone CLI (musl version for Alpine)
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl \
    && chmod +x tailwindcss-linux-x64-musl \
    && mv tailwindcss-linux-x64-musl /usr/local/bin/tailwindcss

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Build Tailwind CSS
RUN tailwindcss -i web/static/css/input.css -o web/static/css/output.css --minify

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o musiq-server .

# Stage 2: Runtime
FROM alpine:latest

# Install FFmpeg for audio/video processing
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/musiq-server .

# Copy static files
COPY --from=builder /app/web/static ./web/static

# Expose port (Koyeb will set PORT env var)
EXPOSE 8080

# Run the server
CMD ["./musiq-server"]
