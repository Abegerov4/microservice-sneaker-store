FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work.sum* ./
COPY notification-service/go.mod notification-service/go.mod
COPY notification-service/go.sum notification-service/go.sum
RUN go work sync 2>/dev/null || true
COPY notification-service/ notification-service/
WORKDIR /workspace/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /notification-service .
CMD ["./notification-service"]
