package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	productpb "sneaker-store/product-service/proto"
	"sneaker-store/order-service/internal/usecase"
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

func (c *ProductGRPCClient) GetProduct(ctx context.Context, id string) (*usecase.ProductInfo, error) {
	resp, err := c.client.GetProduct(ctx, &productpb.GetProductRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return nil, fmt.Errorf("product %s does not exist", id)
		}
		return nil, fmt.Errorf("product service unavailable: %w", err)
	}
	return &usecase.ProductInfo{
		ID:    resp.Id,
		Name:  resp.Name,
		Price: resp.Price,
		Stock: int(resp.Stock),
	}, nil
}
