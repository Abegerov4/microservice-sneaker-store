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

	"sneaker-store/ai-service/internal/ai"
	aicache "sneaker-store/ai-service/internal/cache/redis"
	"sneaker-store/ai-service/internal/client"
	"sneaker-store/ai-service/internal/event"
	chatrepo "sneaker-store/ai-service/internal/repository/postgres"
	grpchandler "sneaker-store/ai-service/internal/transport/grpc"
	"sneaker-store/ai-service/internal/usecase"
	pb "sneaker-store/ai-service/proto"
)

func Run() {
	dbURL := mustEnv("DATABASE_URL")
	geminiKey := mustEnv("OPENAI_API_KEY")
	grpcPort := envOr("GRPC_PORT", "50054")
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
		log.Printf("nats connect: %v — events disabled", err)
		pub = event.NewNATSPublisher(nil)
	} else {
		pub = event.NewNATSPublisher(nc)
		defer nc.Close()
	}

	// Product gRPC client
	productClient, err := client.NewProductGRPCClient(productAddr)
	if err != nil {
		log.Fatalf("product client: %v", err)
	}

	// Gemini client
	geminiClient := ai.NewGeminiClient(geminiKey)

	// Wire
	repo := chatrepo.NewChatRepository(db)
	cache := aicache.NewAICache(redisClient)
	uc := usecase.NewAIUseCase(repo, cache, pub, productClient, geminiClient)
	handler := grpchandler.NewAIHandler(uc)

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterAIServiceServer(srv, handler)

	log.Printf("ai-service listening on :%s", grpcPort)
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
