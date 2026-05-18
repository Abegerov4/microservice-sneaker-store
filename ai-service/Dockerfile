FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work.sum* ./
COPY product-service/go.mod product-service/go.mod
COPY product-service/go.sum product-service/go.sum
COPY ai-service/go.mod ai-service/go.mod
COPY ai-service/go.sum ai-service/go.sum
RUN go work sync 2>/dev/null || true
COPY product-service/ product-service/
COPY ai-service/ ai-service/
WORKDIR /workspace/ai-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /ai-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /ai-service .
COPY --from=builder /workspace/ai-service/migrations ./migrations
EXPOSE 50054
CMD ["./ai-service"]
