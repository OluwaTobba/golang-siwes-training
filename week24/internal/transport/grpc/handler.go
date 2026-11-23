package grpchandler

import (
	"context"
	"errors"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
	"github.com/OluwaTobba/golang-siwes-training/week24/internal/service"
	pb "github.com/OluwaTobba/golang-siwes-training/week24/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	pb.UnimplementedInventoryServiceServer
	svc *service.InventoryService
}

func NewGRPCHandler(svc *service.InventoryService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInsufficientStock):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func toPB(p *domain.Product) *pb.Product {
	return &pb.Product{Id: int32(p.ID), Name: p.Name, Sku: p.SKU, Price: p.Price, Stock: int32(p.Stock)}
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	p := &domain.Product{Name: req.Name, SKU: req.Sku, Price: req.Price, Stock: int(req.Stock)}
	if err := h.svc.CreateProduct(ctx, p); err != nil {
		return nil, toStatus(err)
	}
	return toPB(p), nil
}

func (h *GRPCHandler) AdjustStock(ctx context.Context, req *pb.AdjustStockRequest) (*pb.Product, error) {
	p, err := h.svc.AdjustStock(ctx, int(req.Id), int(req.Delta))
	if err != nil {
		return nil, toStatus(err)
	}
	return toPB(p), nil
}