package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sneaker-store/order-service/internal/model"
	"sneaker-store/order-service/internal/usecase"
	pb "sneaker-store/order-service/proto"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items are required")
	}

	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = usecase.OrderItemInput{
			ProductID: it.ProductId,
			Quantity:  int(it.Quantity),
			Size:      it.Size,
		}
	}

	order, err := h.uc.Create(ctx, usecase.CreateOrderInput{
		UserID:          req.UserId,
		Items:           items,
		ShippingAddress: req.ShippingAddress,
	})
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "create order: %v", err)
	}

	return toProto(order), nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	order, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	return toProto(order), nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, _ *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.uc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list orders: %v", err)
	}
	resp := &pb.ListOrdersResponse{}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return resp, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	order, err := h.uc.UpdateStatus(ctx, req.Id, model.OrderStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "update status: %v", err)
	}
	return toProto(order), nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.OrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	order, err := h.uc.Cancel(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cancel order: %v", err)
	}
	return toProto(order), nil
}

func (h *OrderHandler) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.ListOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	orders, err := h.uc.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get orders by user: %v", err)
	}
	resp := &pb.ListOrdersResponse{}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return resp, nil
}

func (h *OrderHandler) GetOrdersByStatus(ctx context.Context, req *pb.GetOrdersByStatusRequest) (*pb.ListOrdersResponse, error) {
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}
	orders, err := h.uc.GetByStatus(ctx, model.OrderStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "get orders by status: %v", err)
	}
	resp := &pb.ListOrdersResponse{}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return resp, nil
}

func (h *OrderHandler) GetOrderStats(ctx context.Context, _ *pb.GetOrderStatsRequest) (*pb.GetOrderStatsResponse, error) {
	stats, err := h.uc.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get order stats: %v", err)
	}
	return &pb.GetOrderStatsResponse{
		TotalOrders:     int32(stats.TotalOrders),
		PendingOrders:   int32(stats.PendingOrders),
		ConfirmedOrders: int32(stats.ConfirmedOrders),
		ShippedOrders:   int32(stats.ShippedOrders),
		DeliveredOrders: int32(stats.DeliveredOrders),
		CancelledOrders: int32(stats.CancelledOrders),
		TotalRevenue:    stats.TotalRevenue,
	}, nil
}

func (h *OrderHandler) GetOrderItems(ctx context.Context, req *pb.GetOrderItemsRequest) (*pb.GetOrderItemsResponse, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	items, err := h.uc.GetItemsByOrderID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get order items: %v", err)
	}
	resp := &pb.GetOrderItemsResponse{}
	for _, item := range items {
		resp.Items = append(resp.Items, &pb.OrderItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			Price:       item.Price,
			Size:        item.Size,
		})
	}
	return resp, nil
}

func (h *OrderHandler) GetTotalRevenue(ctx context.Context, _ *pb.GetTotalRevenueRequest) (*pb.GetTotalRevenueResponse, error) {
	total, err := h.uc.GetTotalRevenue(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get total revenue: %v", err)
	}
	return &pb.GetTotalRevenueResponse{Total: total}, nil
}

func (h *OrderHandler) GetOrdersByDateRange(ctx context.Context, req *pb.GetOrdersByDateRangeRequest) (*pb.ListOrdersResponse, error) {
	if req.From == "" || req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "from and to are required")
	}
	from, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid from date: %v", err)
	}
	to, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid to date: %v", err)
	}
	orders, err := h.uc.GetByDateRange(ctx, from, to)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get orders by date range: %v", err)
	}
	resp := &pb.ListOrdersResponse{}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return resp, nil
}

func (h *OrderHandler) CountOrdersByUser(ctx context.Context, req *pb.CountOrdersByUserRequest) (*pb.CountOrdersByUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	count, err := h.uc.CountByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "count orders by user: %v", err)
	}
	return &pb.CountOrdersByUserResponse{Count: int32(count)}, nil
}

func toProto(o *model.Order) *pb.OrderResponse {
	resp := &pb.OrderResponse{
		Id:              o.ID,
		UserId:          o.UserID,
		TotalAmount:     o.TotalAmount,
		Status:          string(o.Status),
		ShippingAddress: o.ShippingAddress,
		CreatedAt:       o.CreatedAt.String(),
		UpdatedAt:       o.UpdatedAt.String(),
	}
	for _, item := range o.Items {
		resp.Items = append(resp.Items, &pb.OrderItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			Price:       item.Price,
			Size:        item.Size,
		})
	}
	return resp
}
