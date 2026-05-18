package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"sneaker-store/ai-service/internal/model"
	productpb "sneaker-store/product-service/proto"
)

type ProductGRPCClient struct {
	client productpb.ProductServiceClient
}

func NewProductGRPCClient(addr string) (*ProductGRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to product-service: %w", err)
	}
	return &ProductGRPCClient{client: productpb.NewProductServiceClient(conn)}, nil
}

func (c *ProductGRPCClient) ListProducts(ctx context.Context) ([]*model.ProductInfo, error) {
	resp, err := c.client.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	products := make([]*model.ProductInfo, 0, len(resp.Products))
	for _, p := range resp.Products {
		products = append(products, &model.ProductInfo{
			ID:          p.Id,
			Name:        p.Name,
			Brand:       p.Brand,
			Description: p.Description,
			Price:       p.Price,
			Sizes:       p.Sizes,
			Stock:       p.Stock,
			ImageURL:    p.ImageUrl,
		})
	}
	return products, nil
}
