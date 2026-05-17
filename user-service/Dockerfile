FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.work go.work.sum* ./
COPY user-service/go.mod user-service/go.mod
COPY user-service/go.sum user-service/go.sum
RUN go work sync 2>/dev/null || true
COPY user-service/ user-service/
WORKDIR /workspace/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /user-service .
COPY --from=builder /workspace/user-service/migrations ./migrations
EXPOSE 50053
CMD ["./user-service"]
