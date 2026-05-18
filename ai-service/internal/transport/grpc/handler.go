package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "sneaker-store/ai-service/proto"
	"sneaker-store/ai-service/internal/usecase"
	"sneaker-store/ai-service/internal/model"
)

type AIHandler struct {
	pb.UnimplementedAIServiceServer
	uc *usecase.AIUseCase
}

func NewAIHandler(uc *usecase.AIUseCase) *AIHandler {
	return &AIHandler{uc: uc}
}

func (h *AIHandler) AskSneakerAdvice(ctx context.Context, req *pb.AskSneakerAdviceRequest) (*pb.AskSneakerAdviceResponse, error) {
	if req.Message == "" {
		return nil, status.Error(codes.InvalidArgument, "message is required")
	}

	reply, sessionID, err := h.uc.AskSneakerAdvice(ctx, req.SessionId, req.UserId, req.Message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "advice error: %v", err)
	}

	return &pb.AskSneakerAdviceResponse{
		Reply:     reply,
		SessionId: sessionID,
	}, nil
}

func (h *AIHandler) RecommendSneakers(ctx context.Context, req *pb.RecommendSneakersRequest) (*pb.RecommendSneakersResponse, error) {
	if req.Preferences == "" {
		return nil, status.Error(codes.InvalidArgument, "preferences are required")
	}

	recs, err := h.uc.RecommendSneakers(ctx, req.UserId, req.Preferences, req.Budget, req.Size)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "recommend error: %v", err)
	}

	return &pb.RecommendSneakersResponse{
		Recommendations: toProtoRecs(recs),
	}, nil
}

func (h *AIHandler) SearchSneakersByStyle(ctx context.Context, req *pb.SearchSneakersByStyleRequest) (*pb.SearchSneakersByStyleResponse, error) {
	if req.StyleDescription == "" {
		return nil, status.Error(codes.InvalidArgument, "style_description is required")
	}

	results, summary, err := h.uc.SearchSneakersByStyle(ctx, req.StyleDescription, req.Size, req.MaxPrice)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "style search error: %v", err)
	}

	return &pb.SearchSneakersByStyleResponse{
		Results:   toProtoRecs(results),
		AiSummary: summary,
	}, nil
}

func (h *AIHandler) GetTrendingSneakers(ctx context.Context, req *pb.GetTrendingSneakersRequest) (*pb.GetTrendingSneakersResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	trending, err := h.uc.GetTrendingSneakers(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "trending error: %v", err)
	}

	return &pb.GetTrendingSneakersResponse{
		Sneakers: toProtoTrending(trending),
	}, nil
}

func toProtoRecs(recs []*model.SneakerRecommendation) []*pb.SneakerRecommendation {
	out := make([]*pb.SneakerRecommendation, 0, len(recs))
	for _, r := range recs {
		out = append(out, &pb.SneakerRecommendation{
			ProductId:  r.ProductID,
			Name:       r.Name,
			Brand:      r.Brand,
			Price:      r.Price,
			ImageUrl:   r.ImageURL,
			Reason:     r.Reason,
			MatchScore: r.MatchScore,
		})
	}
	return out
}

func toProtoTrending(trending []*model.TrendingSneaker) []*pb.TrendingSneaker {
	out := make([]*pb.TrendingSneaker, 0, len(trending))
	for _, t := range trending {
		out = append(out, &pb.TrendingSneaker{
			ProductId:   t.ProductID,
			Name:        t.Name,
			Brand:       t.Brand,
			Price:       t.Price,
			ImageUrl:    t.ImageURL,
			TrendReason: t.TrendReason,
			TrendScore:  t.TrendScore,
		})
	}
	return out
}
