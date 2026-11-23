package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/OluwaTobba/golang-siwes-training/week22/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewProductServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Unary call
	p, err := c.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GetProduct(1): %s  price=%.2f  stock=%d\n", p.Name, p.Price, p.Stock)

	// Server-streaming call
	fmt.Println("ListProducts (stream):")
	stream, err := c.ListProducts(ctx, &pb.ListProductsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  #%d %-12s stock=%d\n", item.Id, item.Name, item.Stock)
	}

	// Update, then an error case
	u, err := c.UpdateStock(ctx, &pb.UpdateStockRequest{Id: 2, Delta: -2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UpdateStock(2, -2): stock is now %d\n", u.Stock)

	_, err = c.GetProduct(ctx, &pb.GetProductRequest{Id: 99})
	fmt.Println("GetProduct(99) error:", err)

	_, err = c.UpdateStock(ctx, &pb.UpdateStockRequest{Id: 2, Delta: -100})
	fmt.Println("UpdateStock(2, -100) error:", err)
}
