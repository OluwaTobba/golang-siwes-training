package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "github.com/yourname/siwes/week22/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements ProductServiceServer
type server struct {
	pb.UnimplementedProductServiceServer
	products map[int32]*pb.Product
}

func newServer() *server {
	return &server{products: map[int32]*pb.Product{
		1: {Id: 1, Name: "MacBook Pro", Price: 2499.99, Stock: 10},
		2: {Id: 2, Name: "Dell XPS",    Price: 1899.99, Stock: 5},
		3: {Id: 3, Name: "ThinkPad X1", Price: 1599.99, Stock: 8},
	}}
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	p, ok := s.products[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id)
	}
	return p, nil
}

// Server-side streaming RPC
func (s *server) ListProducts(req *pb.ListProductsRequest, stream pb.ProductService_ListProductsServer) error {
	for _, p := range s.products {
		if err := stream.Send(p); err != nil { return err }
	}
	return nil
}

func (s *server) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.Product, error) {
	p, ok := s.products[req.Id]
	if !ok { return nil, status.Errorf(codes.NotFound, "product %d not found", req.Id) }
	p.Stock += req.Delta
	if p.Stock < 0 { return nil, status.Error(codes.InvalidArgument, "stock cannot be negative") }
	return p, nil
}

// Logging interceptor
func loggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("gRPC call: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil { log.Printf("error: %v", err) }
	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil { log.Fatal(err) }
	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	pb.RegisterProductServiceServer(s, newServer())
	fmt.Println("gRPC server on :50051")
	log.Fatal(s.Serve(lis))
}