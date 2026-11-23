package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"sync"

	pb "github.com/OluwaTobba/golang-siwes-training/week22/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server implements ProductServiceServer
type server struct {
	pb.UnimplementedProductServiceServer
	mu       sync.Mutex
	products map[int32]*pb.Product
}

func newServer() *server {
	return &server{products: map[int32]*pb.Product{
		1: {Id: 1, Name: "MacBook Pro", Price: 2499.99, Stock: 10},
		2: {Id: 2, Name: "Dell XPS", Price: 1899.99, Stock: 5},
		3: {Id: 3, Name: "ThinkPad X1", Price: 1599.99, Stock: 8},
	}}
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	return p, nil
}

// Server-side streaming RPC
func (s *server) ListProducts(req *pb.ListProductsRequest, stream pb.ProductService_ListProductsServer) error {
	s.mu.Lock()
	ids := make([]int, 0, len(s.products))
	for id := range s.products {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	list := make([]*pb.Product, 0, len(ids))
	for _, id := range ids {
		list = append(list, s.products[int32(id)])
	}
	s.mu.Unlock()

	for _, p := range list {
		if err := stream.Send(p); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	if p.Stock+req.Delta < 0 {
		return nil, status.Error(codes.InvalidArgument, "stock cannot be negative")
	}
	p.Stock += req.Delta
	return p, nil
}

// Logging interceptor
func loggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("gRPC call: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	pb.RegisterProductServiceServer(s, newServer())
	fmt.Println("gRPC server on :50051")
	log.Fatal(s.Serve(lis))
}
