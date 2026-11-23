package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/OluwaTobba/golang-siwes-training/week24/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewInventoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sku := fmt.Sprintf("GRPC-%d", time.Now().Unix())
	p, err := c.CreateProduct(ctx, &pb.CreateProductRequest{Name: "Keyboard", Sku: sku, Price: 49.99, Stock: 20})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CreateProduct: id=%d sku=%s stock=%d\n", p.Id, p.Sku, p.Stock)

	u, err := c.AdjustStock(ctx, &pb.AdjustStockRequest{Id: p.Id, Delta: -5})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AdjustStock(-5): stock is now %d\n", u.Stock)

	_, err = c.AdjustStock(ctx, &pb.AdjustStockRequest{Id: p.Id, Delta: -1000})
	fmt.Println("AdjustStock(-1000) error:", err)
}