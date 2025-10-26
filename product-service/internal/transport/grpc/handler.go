package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sneaker-store/product-service/internal/model"
	"sneaker-store/product-service/internal/usecase"
	pb "sneaker-store/product-service/proto"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	uc *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	p := &model.Product{
		Name:        req.Name,
		Brand:       req.Brand,
		Description: req.Description,
		Price:       req.Price,
		Sizes:       req.Sizes,
		Stock:       int(req.Stock),
		ImageURL:    req.ImageUrl,
	}
	created, err := h.uc.Create(ctx, p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create product: %v", err)
	}
	return toProto(created), nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	p, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}
	return toProto(p), nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, _ *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.uc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list products: %v", err)
	}
	resp := &pb.ListProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProto(p))
	}
	return resp, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	p := &model.Product{
		ID:          req.Id,
		Name:        req.Name,
		Brand:       req.Brand,
		Description: req.Description,
		Price:       req.Price,
		Sizes:       req.Sizes,
		Stock:       int(req.Stock),
		ImageURL:    req.ImageUrl,
	}
	updated, err := h.uc.Update(ctx, p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update product: %v", err)
	}
	return toProto(updated), nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "delete product: %v", err)
	}
	return &pb.DeleteProductResponse{Success: true}, nil
}

func (h *ProductHandler) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.uc.Search(ctx, req.Brand, req.MinPrice, req.MaxPrice, req.Size)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search products: %v", err)
	}
	resp := &pb.ListProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProto(p))
	}
	return resp, nil
}

func (h *ProductHandler) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	newStock, err := h.uc.UpdateStock(ctx, req.Id, int(req.Delta))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update stock: %v", err)
	}
	return &pb.UpdateStockResponse{Id: req.Id, NewStock: int32(newStock)}, nil
}

func (h *ProductHandler) GetProductsByBrand(ctx context.Context, req *pb.GetProductsByBrandRequest) (*pb.ListProductsResponse, error) {
	if req.Brand == "" {
		return nil, status.Error(codes.InvalidArgument, "brand is required")
	}
	products, err := h.uc.GetByBrand(ctx, req.Brand)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get products by brand: %v", err)
	}
	resp := &pb.ListProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProto(p))
	}
	return resp, nil
}

func (h *ProductHandler) GetLowStockProducts(ctx context.Context, req *pb.GetLowStockRequest) (*pb.ListProductsResponse, error) {
	threshold := int(req.Threshold)
	if threshold <= 0 {
		threshold = 10
	}
	products, err := h.uc.GetLowStock(ctx, threshold)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get low stock: %v", err)
	}
	resp := &pb.ListProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProto(p))
	}
	return resp, nil
}

func (h *ProductHandler) GetBrands(ctx context.Context, _ *pb.GetBrandsRequest) (*pb.GetBrandsResponse, error) {
	brands, err := h.uc.GetBrands(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get brands: %v", err)
	}
	return &pb.GetBrandsResponse{Brands: brands}, nil
}

func (h *ProductHandler) GetProductStats(ctx context.Context, _ *pb.GetProductStatsRequest) (*pb.GetProductStatsResponse, error) {
	stats, err := h.uc.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get stats: %v", err)
	}
	return &pb.GetProductStatsResponse{
		TotalProducts: int32(stats.TotalProducts),
		TotalBrands:   int32(stats.TotalBrands),
		TotalStock:    int32(stats.TotalStock),
		AveragePrice:  stats.AveragePrice,
	}, nil
}

func (h *ProductHandler) BulkDeleteProducts(ctx context.Context, req *pb.BulkDeleteProductsRequest) (*pb.BulkDeleteProductsResponse, error) {
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ids are required")
	}
	count, err := h.uc.BulkDelete(ctx, req.Ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "bulk delete: %v", err)
	}
	return &pb.BulkDeleteProductsResponse{DeletedCount: int32(count)}, nil
}

func toProto(p *model.Product) *pb.ProductResponse {
	return &pb.ProductResponse{
		Id:          p.ID,
		Name:        p.Name,
		Brand:       p.Brand,
		Description: p.Description,
		Price:       p.Price,
		Sizes:       p.Sizes,
		Stock:       int32(p.Stock),
		ImageUrl:    p.ImageURL,
		CreatedAt:   p.CreatedAt.String(),
		UpdatedAt:   p.UpdatedAt.String(),
	}
}
