package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	ordercache "sneaker-store/order-service/internal/cache/redis"
	"sneaker-store/order-service/internal/client"
	"sneaker-store/order-service/internal/event"
	orderrepo "sneaker-store/order-service/internal/repository/postgres"
	grpchandler "sneaker-store/order-service/internal/transport/grpc"
	"sneaker-store/order-service/internal/usecase"
	pb "sneaker-store/order-service/proto"
)

func Run() {
	dbURL := mustEnv("DATABASE_URL")
	grpcPort := envOr("GRPC_PORT", "50052")
	natsURL := envOr("NATS_URL", nats.DefaultURL)
	redisURL := envOr("REDIS_URL", "redis://localhost:6379")
	productAddr := envOr("PRODUCT_SERVICE_ADDR", "localhost:50051")

	// Migrations
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("migrate create: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations applied")

	// DB pool
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("pgx pool: %v", err)
	}
	defer db.Close()

	// Redis
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("redis parse url: %v — caching disabled", err)
	}
	redisClient := redis.NewClient(opt)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("redis unavailable: %v — caching best-effort", err)
	}

	// NATS
	var pub *event.NATSPublisher
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("nats connect failed: %v — events disabled", err)
		pub = event.NewNATSPublisherConn(nil)
	} else {
		pub = event.NewNATSPublisherConn(nc)
		defer nc.Close()
	}

	// Product gRPC client
	productClient, err := client.NewProductGRPCClient(productAddr)
	if err != nil {
		log.Fatalf("product client: %v", err)
	}

	// Wire
	repo := orderrepo.NewOrderRepository(db)
	cache := ordercache.NewOrderCache(redisClient)
	uc := usecase.NewOrderUseCase(repo, cache, pub, productClient)
	handler := grpchandler.NewOrderHandler(uc)

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterOrderServiceServer(srv, handler)

	log.Printf("order-service listening on :%s", grpcPort)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
