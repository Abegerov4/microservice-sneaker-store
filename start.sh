#!/usr/bin/env bash
set -e

# ── Colors ───────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; NC='\033[0m'; BOLD='\033[1m'

info()    { echo -e "${CYAN}[INFO]${NC}  $1"; }
success() { echo -e "${GREEN}[OK]${NC}    $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}  $1"; }
header()  { echo -e "\n${BOLD}$1${NC}"; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_DIR="$SCRIPT_DIR/logs"
mkdir -p "$LOG_DIR"

# ── 1. Infrastructure ─────────────────────────────────────────────────────────
header "1. Starting infrastructure (Docker)"

# PostgreSQL
if docker inspect postgres &>/dev/null; then
  docker start postgres &>/dev/null && success "postgres started"
else
  info "Creating postgres container..."
  docker run -d --name postgres \
    -e POSTGRES_PASSWORD=postgres \
    -p 5432:5432 \
    postgres:16 &>/dev/null
  success "postgres created"
fi

# NATS
if docker inspect nats &>/dev/null; then
  docker start nats &>/dev/null && success "nats started"
else
  info "Creating nats container..."
  docker run -d --name nats -p 4222:4222 nats:latest &>/dev/null
  success "nats created"
fi

# Redis
if docker inspect redis &>/dev/null; then
  docker start redis &>/dev/null && success "redis started"
else
  info "Creating redis container..."
  docker run -d --name redis -p 6379:6379 redis:7-alpine &>/dev/null
  success "redis created"
fi

# ── 2. Wait for Postgres ──────────────────────────────────────────────────────
header "2. Waiting for PostgreSQL to be ready"
for i in {1..20}; do
  if docker exec postgres pg_isready -U postgres &>/dev/null; then
    success "PostgreSQL is ready"
    break
  fi
  if [ $i -eq 20 ]; then
    echo -e "${RED}[ERROR]${NC} PostgreSQL did not start in time"; exit 1
  fi
  sleep 1
done

# ── 3. Create databases ───────────────────────────────────────────────────────
header "3. Creating databases"
PG_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
for db in sneaker_products sneaker_orders sneaker_users sneaker_ai; do
  exists=$(psql "$PG_URL" -tAc "SELECT 1 FROM pg_database WHERE datname='$db'" 2>/dev/null || true)
  if [ "$exists" = "1" ]; then
    success "$db already exists"
  else
    psql "$PG_URL" -c "CREATE DATABASE $db;" &>/dev/null && success "$db created" || warn "$db creation failed"
  fi
done

# ── 4. Go services ────────────────────────────────────────────────────────────
header "4. Starting Go services"

PIDS=()

start_service() {
  local name="$1"; local log="$LOG_DIR/${name}.log"
  shift 1
  info "Starting $name..."
  env "$@" go run ./cmd > "$log" 2>&1 &
  local pid=$!
  PIDS+=("$pid")
  success "$name → PID $pid (log: logs/${name}.log)"
}

cd "$SCRIPT_DIR/product-service"
start_service "product-service" \
  "DATABASE_URL=postgres://postgres:postgres@localhost:5432/sneaker_products?sslmode=disable" \
  "NATS_URL=nats://localhost:4222" \
  "REDIS_URL=redis://localhost:6379" \
  "GRPC_PORT=50051"

cd "$SCRIPT_DIR/user-service"
start_service "user-service" \
  "DATABASE_URL=postgres://postgres:postgres@localhost:5432/sneaker_users?sslmode=disable" \
  "NATS_URL=nats://localhost:4222" \
  "REDIS_URL=redis://localhost:6379" \
  "GRPC_PORT=50053" \
  "ADMIN_EMAIL=${ADMIN_EMAIL:-admin@sneakerstore.com}" \
  "ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin123}"

# Wait a moment for product-service to bind its port before order-service connects
sleep 3

cd "$SCRIPT_DIR/order-service"
start_service "order-service" \
  "DATABASE_URL=postgres://postgres:postgres@localhost:5432/sneaker_orders?sslmode=disable" \
  "NATS_URL=nats://localhost:4222" \
  "REDIS_URL=redis://localhost:6379" \
  "GRPC_PORT=50052" \
  "PRODUCT_SERVICE_ADDR=localhost:50051"

cd "$SCRIPT_DIR/notification-service"
start_service "notification-service" \
  "NATS_URL=nats://localhost:4222" \
  "SMTP_USERNAME=${SMTP_USERNAME:-aldiar12378@gmail.com}" \
  "SMTP_PASSWORD=${SMTP_PASSWORD:-${SMTP_PASSWORD}}" \
  "SMTP_FROM=${SMTP_FROM:-aldiar12378@gmail.com}" \
  "NOTIFY_EMAIL=${NOTIFY_EMAIL:-aldiar12378@gmail.com}"

cd "$SCRIPT_DIR/ai-service"
start_service "ai-service" \
  "DATABASE_URL=postgres://postgres:postgres@localhost:5432/sneaker_ai?sslmode=disable" \
  "NATS_URL=nats://localhost:4222" \
  "REDIS_URL=redis://localhost:6379" \
  "GRPC_PORT=50054" \
  "PRODUCT_SERVICE_ADDR=localhost:50051" \
  "OPENAI_API_KEY=${OPENAI_API_KEY:-${OPENAI_API_KEY}}"

sleep 2

cd "$SCRIPT_DIR/api-gateway"
start_service "api-gateway" \
  "PRODUCT_SERVICE_ADDR=localhost:50051" \
  "ORDER_SERVICE_ADDR=localhost:50052" \
  "USER_SERVICE_ADDR=localhost:50053" \
  "AI_SERVICE_ADDR=localhost:50054" \
  "HTTP_PORT=8080" \
  "JWT_SECRET=${JWT_SECRET:-sneaker-store-jwt-secret-2026}"

# ── 5. Frontend ───────────────────────────────────────────────────────────────
header "5. Starting Next.js frontend"
cd "$SCRIPT_DIR/frontend"
info "Starting frontend (npm run dev)..."
npm run dev > "$LOG_DIR/frontend.log" 2>&1 &
FRONTEND_PID=$!
PIDS+=("$FRONTEND_PID")
success "frontend → PID $FRONTEND_PID (log: logs/frontend.log)"

# ── 6. Done ───────────────────────────────────────────────────────────────────
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}${BOLD}  All services started!${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "  🌐 Frontend:    ${CYAN}http://localhost:3000${NC}"
echo -e "  🔐 Admin Panel: ${CYAN}http://localhost:3000/admin${NC}"
echo -e "  🔧 API Gateway: ${CYAN}http://localhost:8080${NC}"
echo -e "  🤖 AI Service:  ${CYAN}gRPC :50054${NC}"
echo -e "  📋 Logs:        ${CYAN}./logs/${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Press ${BOLD}Ctrl+C${NC} to stop all services."

# ── Cleanup on Ctrl+C ─────────────────────────────────────────────────────────
cleanup() {
  echo ""
  warn "Stopping all services..."
  for pid in "${PIDS[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  success "Done."
}
trap cleanup INT TERM

wait
