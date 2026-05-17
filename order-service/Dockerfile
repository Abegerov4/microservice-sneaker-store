FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work.sum* ./
COPY product-service/go.mod product-service/go.mod
COPY product-service/go.sum product-service/go.sum
COPY order-service/go.mod order-service/go.mod
COPY order-service/go.sum order-service/go.sum
RUN go work sync 2>/dev/null || true
COPY product-service/ product-service/
COPY order-service/ order-service/
WORKDIR /workspace/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /order-service .
COPY --from=builder /workspace/order-service/migrations ./migrations
EXPOSE 50052
CMD ["./order-service"]
