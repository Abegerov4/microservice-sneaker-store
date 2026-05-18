FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work.sum* ./
COPY product-service/go.mod product-service/go.mod
COPY product-service/go.sum product-service/go.sum
COPY order-service/go.mod order-service/go.mod
COPY order-service/go.sum order-service/go.sum
COPY user-service/go.mod user-service/go.mod
COPY user-service/go.sum user-service/go.sum
COPY ai-service/go.mod ai-service/go.mod
COPY ai-service/go.sum ai-service/go.sum
COPY api-gateway/go.mod api-gateway/go.mod
COPY api-gateway/go.sum api-gateway/go.sum
RUN go work sync 2>/dev/null || true
COPY product-service/ product-service/
COPY order-service/ order-service/
COPY user-service/ user-service/
COPY ai-service/ ai-service/
COPY api-gateway/ api-gateway/
WORKDIR /workspace/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
