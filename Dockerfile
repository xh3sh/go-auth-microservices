# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build args to specify which service to build
ARG SERVICE_PATH

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service ./cmd/${SERVICE_PATH}

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/service .

# Copy static assets and views if they exist (needed for frontend)
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views
COPY --from=builder /app/.env ./.env

EXPOSE 8080 8081 8082 8083 8084

CMD ["./service"]
