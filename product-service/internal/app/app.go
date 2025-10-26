package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	productcache "sneaker-store/product-service/internal/cache/redis"
	"sneaker-store/product-service/internal/event"
	productrepo "sneaker-store/product-service/internal/repository/postgres"
	grpchandler "sneaker-store/product-service/internal/transport/grpc"
	"sneaker-store/product-service/internal/usecase"
	pb "sneaker-store/product-service/proto"
)

func Run() {
	dbURL := mustEnv("DATABASE_URL")
	grpcPort := envOr("GRPC_PORT", "50051")
	natsURL := envOr("NATS_URL", nats.DefaultURL)
	redisURL := envOr("REDIS_URL", "redis://localhost:6379")

	// Run migrations
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
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("nats connect failed: %v — events disabled", err)
	}
	var pub *event.NATSPublisher
	if nc != nil {
		pub = event.NewNATSPublisherConn(nc)
		defer nc.Close()
	} else {
		pub = &event.NATSPublisher{}
	}

	// Wire
	repo := productrepo.NewProductRepository(db)
	cache := productcache.NewProductCache(redisClient)
	uc := usecase.NewProductUseCase(repo, cache, pub)
	handler := grpchandler.NewProductHandler(uc)

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterProductServiceServer(srv, handler)

	log.Printf("product-service listening on :%s", grpcPort)
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

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
