FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work.sum* ./
COPY product-service/go.mod product-service/go.mod
COPY product-service/go.sum product-service/go.sum
RUN go work sync 2>/dev/null || true
COPY product-service/ product-service/
WORKDIR /workspace/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /product-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /product-service .
COPY --from=builder /workspace/product-service/migrations ./migrations
EXPOSE 50051
CMD ["./product-service"]
