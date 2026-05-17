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

	usercache "sneaker-store/user-service/internal/cache/redis"
	"sneaker-store/user-service/internal/event"
	userrepo "sneaker-store/user-service/internal/repository/postgres"
	grpchandler "sneaker-store/user-service/internal/transport/grpc"
	"sneaker-store/user-service/internal/usecase"
	pb "sneaker-store/user-service/proto"
)

func Run() {
	dbURL := mustEnv("DATABASE_URL")
	grpcPort := envOr("GRPC_PORT", "50053")
	natsURL := envOr("NATS_URL", nats.DefaultURL)
	redisURL := envOr("REDIS_URL", "redis://localhost:6379")
	adminEmail := envOr("ADMIN_EMAIL", "admin@sneakerstore.com")
	adminPassword := envOr("ADMIN_PASSWORD", "admin123")

	// Migrations
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("migrate create: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations applied")

	// DB
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

	// Wire
	repo := userrepo.NewUserRepository(db)
	cache := usercache.NewUserCache(redisClient)
	uc := usecase.NewUserUseCase(repo, cache, pub)

	// Seed admin account
	if err := uc.EnsureAdmin(context.Background(), adminEmail, adminPassword, "Administrator"); err != nil {
		log.Printf("ensure admin: %v", err)
	} else {
		log.Printf("admin account ready: %s", adminEmail)
	}

	handler := grpchandler.NewUserHandler(uc)

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, handler)

	log.Printf("user-service listening on :%s", grpcPort)
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
