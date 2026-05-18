# Sneaker Store — Final Project Backend

A microservices-based online sneaker store built with Go, gRPC, PostgreSQL, Redis, and NATS.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     API Gateway :8080                    │
│              (HTTP REST → gRPC routing)                  │
└──────────┬────────────────┬────────────────┬────────────┘
           │                │                │
           ▼                ▼                ▼
  ┌────────────────┐ ┌────────────────┐ ┌────────────────┐
  │ Product Service│ │  Order Service │ │  User Service  │
  │    :50051      │ │    :50052      │ │    :50053      │
  │  PostgreSQL    │ │  PostgreSQL    │ │  PostgreSQL    │
  │  Redis cache   │ │  Redis cache   │ │  Redis cache   │
  └───────┬────────┘ └────────┬───────┘ └───────┬────────┘
          │  NATS publish     │ NATS publish      │ NATS publish
          ▼                   ▼                   ▼
  ┌────────────────────────────────────────────────────────┐
  │               NATS Message Broker :4222                │
  └────────────────────────────┬───────────────────────────┘
                               │
                               ▼
              ┌────────────────────────────────┐
              │      Notification Service      │
              │  (subscriber + Gmail SMTP)     │
              └────────────────────────────────┘
```

## Services

| Service              | Port  | Description                                    |
|----------------------|-------|------------------------------------------------|
| product-service      | 50051 | Sneaker catalog management                     |
| order-service        | 50052 | Order management (calls product-service gRPC)  |
| user-service         | 50053 | User accounts and authentication               |
| notification-service | —     | NATS subscriber, sends Gmail emails            |
| api-gateway          | 8080  | HTTP REST → gRPC gateway                       |

## gRPC Endpoints (36 total)

### Product Service (12)
| RPC                  | Description                                   |
|----------------------|-----------------------------------------------|
| CreateProduct        | Add a new sneaker                             |
| GetProduct           | Get sneaker by ID                             |
| ListProducts         | List all sneakers                             |
| UpdateProduct        | Update sneaker info                           |
| DeleteProduct        | Remove sneaker                                |
| SearchProducts       | Filter by brand, price range, size            |
| UpdateStock          | Adjust stock quantity by delta                |
| GetProductsByBrand   | List all sneakers of a specific brand         |
| GetLowStockProducts  | List sneakers with stock below threshold      |
| GetBrands            | List all distinct brand names                 |
| GetProductStats      | Aggregate stats (count, avg price, stock)     |
| BulkDeleteProducts   | Delete multiple sneakers by ID list           |

### Order Service (12)
| RPC                    | Description                                 |
|------------------------|---------------------------------------------|
| CreateOrder            | Place a new order (validates products)      |
| GetOrder               | Get order by ID                             |
| ListOrders             | List all orders                             |
| UpdateOrderStatus      | Change order status                         |
| CancelOrder            | Cancel an order                             |
| GetOrdersByUser        | Get all orders for a user                   |
| GetOrdersByStatus      | Filter orders by status                     |
| GetOrderStats          | Count orders by status + total revenue      |
| GetOrderItems          | Return line items for a specific order      |
| GetTotalRevenue        | Sum of all non-cancelled order totals       |
| GetOrdersByDateRange   | Filter orders by created_at range           |
| CountOrdersByUser      | Count total orders for a user               |

### User Service (12)
| RPC                | Description                                  |
|--------------------|----------------------------------------------|
| CreateUser         | Register a new user (bcrypt password)        |
| GetUser            | Get user by ID                               |
| UpdateUser         | Update name / phone                          |
| DeleteUser         | Delete user account                          |
| AuthenticateUser   | Login with email + password                  |
| GetUserByEmail     | Look up user by email address                |
| ChangePassword     | Change password with old + new               |
| ListUsers          | Paginated list of all users                  |
| SearchUsers        | Search users by name or email                |
| GetUserStats       | Total and active user counts                 |
| UpdateUserStatus   | Activate or deactivate a user account        |
| ResetPassword      | Admin password reset (no old pass required)  |

## Prerequisites

```bash
# PostgreSQL
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16

# Create separate databases
psql -U postgres -c "CREATE DATABASE sneaker_products;"
psql -U postgres -c "CREATE DATABASE sneaker_orders;"
psql -U postgres -c "CREATE DATABASE sneaker_users;"

# NATS
docker run -d --name nats -p 4222:4222 nats:latest

# Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

## How to regenerate proto stubs

```bash
export PATH="$PATH:$(go env GOPATH)/bin"

# Product
cd product-service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/product.proto

# Order
cd ../order-service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/order.proto

# User
cd ../user-service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto
```

## Running the Services

Start in this order (dependencies first):

```bash
# 1. Product Service
cd product-service
DATABASE_URL="postgres://postgres:postgres@localhost:5432/sneaker_products?sslmode=disable" \
NATS_URL="nats://localhost:4222" \
REDIS_URL="redis://localhost:6379" \
GRPC_PORT=50051 \
go run ./cmd

# 2. User Service
cd user-service
DATABASE_URL="postgres://postgres:postgres@localhost:5432/sneaker_users?sslmode=disable" \
NATS_URL="nats://localhost:4222" \
REDIS_URL="redis://localhost:6379" \
GRPC_PORT=50053 \
go run ./cmd

# 3. Order Service (needs product-service running)
cd order-service
DATABASE_URL="postgres://postgres:postgres@localhost:5432/sneaker_orders?sslmode=disable" \
NATS_URL="nats://localhost:4222" \
REDIS_URL="redis://localhost:6379" \
GRPC_PORT=50052 \
PRODUCT_SERVICE_ADDR=localhost:50051 \
go run ./cmd

# 4. Notification Service
cd notification-service
NATS_URL="nats://localhost:4222" \
SMTP_USERNAME="your@gmail.com" \
SMTP_PASSWORD="your-app-password" \
SMTP_FROM="your@gmail.com" \
NOTIFY_EMAIL="notifications@example.com" \
go run ./cmd

# 5. API Gateway
cd api-gateway
PRODUCT_SERVICE_ADDR=localhost:50051 \
ORDER_SERVICE_ADDR=localhost:50052 \
USER_SERVICE_ADDR=localhost:50053 \
HTTP_PORT=8080 \
go run ./cmd
```

## Environment Variables

| Variable               | Service              | Default              | Description                    |
|------------------------|----------------------|----------------------|--------------------------------|
| DATABASE_URL           | product/order/user   | required             | PostgreSQL connection string   |
| GRPC_PORT              | product/order/user   | 50051/50052/50053    | gRPC listen port               |
| NATS_URL               | all                  | nats://localhost:4222| NATS connection URL            |
| REDIS_URL              | product/order/user   | redis://localhost:6379| Redis connection URL           |
| PRODUCT_SERVICE_ADDR   | order, api-gateway   | localhost:50051      | Product service gRPC address   |
| ORDER_SERVICE_ADDR     | api-gateway          | localhost:50052      | Order service gRPC address     |
| USER_SERVICE_ADDR      | api-gateway          | localhost:50053      | User service gRPC address      |
| SMTP_HOST              | notification         | smtp.gmail.com       | SMTP server                    |
| SMTP_PORT              | notification         | 587                  | SMTP port                      |
| SMTP_USERNAME          | notification         | —                    | Gmail address                  |
| SMTP_PASSWORD          | notification         | —                    | Gmail App Password             |
| HTTP_PORT              | api-gateway          | 8080                 | HTTP listen port               |

## REST API Examples (via API Gateway)

```bash
# Create a product
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Air Max 90","brand":"Nike","price":120.00,"sizes":["40","41","42","43"],"stock":50}'

# List products
curl http://localhost:8080/api/v1/products

# Search products
curl "http://localhost:8080/api/v1/products/search?brand=Nike"

# Register user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123","full_name":"John Doe"}'

# Login
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}'

# Create order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":"<user-id>","items":[{"product_id":"<product-id>","quantity":1,"size":"42"}],"shipping_address":"Almaty, Kazakhstan"}'

# Update order status
curl -X PATCH http://localhost:8080/api/v1/orders/<order-id>/status \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed"}'
```

## NATS Events Published

| Subject                  | Trigger                          | Payload Fields                              |
|--------------------------|----------------------------------|---------------------------------------------|
| products.created         | CreateProduct succeeds           | event_type, occurred_at, id, name, brand    |
| products.updated         | UpdateProduct succeeds           | event_type, occurred_at, id, name, brand    |
| orders.created           | CreateOrder succeeds             | event_type, occurred_at, id, user_id, status|
| orders.status_updated    | UpdateOrderStatus succeeds       | event_type, occurred_at, id, old/new_status |
| users.registered         | CreateUser succeeds              | event_type, occurred_at, id, email          |

## Broker Choice: NATS

NATS (Core) was chosen for:
- Simpler setup — single binary, no broker configuration
- Stateless pub/sub suitable for notifications
- Lower latency than RabbitMQ for fire-and-forget patterns
- Preferred by the course rubric

**Trade-off vs RabbitMQ:** NATS Core is fire-and-forget — if the notification service is down when an event is published, the event is lost. RabbitMQ fanout queues persist messages. In production, NATS JetStream or the Outbox pattern would address this.

## Clean Architecture

Each service follows the same layered structure:
```
internal/
  model/       — pure domain structs, no framework imports
  usecase/     — business logic + interfaces (no infrastructure)
  repository/postgres/ — PostgreSQL implementation of repo interface
  cache/redis/ — Redis implementation of cache interface
  event/       — NATS event publisher implementation
  transport/grpc/ — thin gRPC handler: parse proto → call usecase → return proto
  app/         — wiring: connects all layers, runs migrations, starts server
```

**Dependency rule:** domain ← usecase ← infrastructure. Use cases depend on interfaces, never on concrete Redis/PostgreSQL/NATS types.

## Database Migrations

Migrations run automatically on service startup. To run manually:

```bash
# Install golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/sneaker_products?sslmode=disable" up

# Rollback
migrate -path ./migrations -database "..." down 1
```

## Failure Scenarios

| Situation                     | Behaviour                                                  |
|-------------------------------|------------------------------------------------------------|
| Product service down          | Order creation returns codes.FailedPrecondition            |
| DB unavailable on startup     | Service exits with non-zero code                           |
| Redis unavailable             | Caching skipped, all reads go to DB (best-effort)         |
| NATS unavailable              | Service starts normally, events not published (logged)     |
| SMTP not configured           | Email logged but not sent (graceful degradation)           |
