FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY backend/go.mod ./
COPY backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN go build -o /app/bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/server /app/server
COPY backend/migrations /app/migrations

ENV PORT=8080
ENV MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

CMD ["/app/server"]
