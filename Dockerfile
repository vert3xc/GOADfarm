# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git (needed for go get) and build deps
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy static files and templates
COPY --from=builder /app/server/frontend /app/server/frontend
COPY --from=builder /app/server/templates /app/server/templates

# Copy built binary
COPY --from=builder /app/server/server /app/server

# Copy any other needed files (e.g., migrations, config)
# COPY --from=builder /app/server/config.yaml /app/server/

EXPOSE 5001

CMD ["/app/server/server"]