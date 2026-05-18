package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	aipb "sneaker-store/ai-service/proto"
	orderpb "sneaker-store/order-service/proto"
	productpb "sneaker-store/product-service/proto"
	userpb "sneaker-store/user-service/proto"
)

type Clients struct {
	Product productpb.ProductServiceClient
	Order   orderpb.OrderServiceClient
	User    userpb.UserServiceClient
	AI      aipb.AIServiceClient
}

func NewClients(productAddr, orderAddr, userAddr, aiAddr string) (*Clients, error) {
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("product client: %w", err)
	}

	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("order client: %w", err)
	}

	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("user client: %w", err)
	}

	aiConn, err := grpc.NewClient(aiAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("ai client: %w", err)
	}

	return &Clients{
		Product: productpb.NewProductServiceClient(productConn),
		Order:   orderpb.NewOrderServiceClient(orderConn),
		User:    userpb.NewUserServiceClient(userConn),
		AI:      aipb.NewAIServiceClient(aiConn),
	}, nil
}
