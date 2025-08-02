FROM golang:1.24.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server_bin ./server/main.go

# ---- Stage 2: Run ----
FROM debian:bullseye-slim

COPY --from=builder /app/server_bin /server/server

COPY --from=builder /app/server/frontend/static /server/frontend/static
COPY --from=builder /app/server/frontend/templates /server/frontend/templates

WORKDIR /server

RUN mkdir uploads && chmod 777 uploads

CMD ["./server"]