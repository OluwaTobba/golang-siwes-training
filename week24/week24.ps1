$mod = "github.com/OluwaTobba/golang-siwes-training/week24"

# 0. Tools: protoc + Go plugins
$env:Path = [Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [Environment]::GetEnvironmentVariable("Path","User")
$env:Path += ";" + (Join-Path (go env GOPATH) "bin")
if (-not (Get-Command protoc -ErrorAction SilentlyContinue)) { Write-Host "protoc not found - install it as in Week 22 first."; return }
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 1. Folders and module
New-Item -ItemType Directory -Force proto, internal/domain, internal/service, internal/repository, internal/transport/http, internal/transport/grpc, cmd/server, cmd/client, k8s | Out-Null
if (-not (Test-Path go.mod)) { go mod init $mod }

# 2. Proto
@'
syntax = "proto3";

option go_package = "github.com/OluwaTobba/golang-siwes-training/week24/proto;pb";

package inventory;

message Product {
  int32 id = 1;
  string name = 2;
  string sku = 3;
  double price = 4;
  int32 stock = 5;
}

message CreateProductRequest {
  string name = 1;
  string sku = 2;
  double price = 3;
  int32 stock = 4;
}

message AdjustStockRequest {
  int32 id = 1;
  int32 delta = 2;
}

service InventoryService {
  rpc CreateProduct (CreateProductRequest) returns (Product);
  rpc AdjustStock (AdjustStockRequest) returns (Product);
}
'@ | Set-Content -Encoding ascii proto/inventory.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/inventory.proto

# 3. Domain
@'
package domain

import (
	"errors"
	"fmt"
	"time"
)

// Domain errors
var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidInput      = errors.New("invalid input")
)

type Product struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	SKU        string    `json:"sku" gorm:"uniqueIndex;not null"`
	CategoryID int       `json:"category_id"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if p.SKU == "" {
		return fmt.Errorf("%w: sku is required", ErrInvalidInput)
	}
	if p.Price < 0 {
		return fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}
	if p.Stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	}
	return nil
}

func (p *Product) AdjustStock(delta int) error {
	newStock := p.Stock + delta
	if newStock < 0 {
		return fmt.Errorf("%w: have %d, need %d", ErrInsufficientStock, p.Stock, -delta)
	}
	p.Stock = newStock
	return nil
}
'@ | Set-Content -Encoding ascii internal/domain/product.go

# 4. Service
@'
package service

import (
	"context"
	"fmt"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
)

type ProductRepo interface {
	Create(ctx context.Context, p *domain.Product) error
	GetByID(ctx context.Context, id int) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id int) error
}

type InventoryService struct{ repo ProductRepo }

func NewInventoryService(repo ProductRepo) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) CreateProduct(ctx context.Context, p *domain.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, p)
}

func (s *InventoryService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *InventoryService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *InventoryService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *InventoryService) AdjustStock(ctx context.Context, id, delta int) (*domain.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.AdjustStock(delta); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("updating stock: %w", err)
	}
	return p, nil
}
'@ | Set-Content -Encoding ascii internal/service/inventory.go

# 5. Repository (GORM)
@'
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
	"gorm.io/gorm"
)

type GORMProductRepo struct{ db *gorm.DB }

func NewGORMProductRepo(db *gorm.DB) *GORMProductRepo {
	return &GORMProductRepo{db: db}
}

func (r *GORMProductRepo) Create(ctx context.Context, p *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("%w: sku %q already exists", domain.ErrInvalidInput, p.SKU)
		}
		return err
	}
	return nil
}

func (r *GORMProductRepo) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var p domain.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *GORMProductRepo) List(ctx context.Context) ([]*domain.Product, error) {
	var ps []*domain.Product
	if err := r.db.WithContext(ctx).Order("id").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (r *GORMProductRepo) Update(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *GORMProductRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&domain.Product{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}
'@ | Set-Content -Encoding ascii internal/repository/gorm_repo.go

# 6. REST handler
@'
package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
	"github.com/OluwaTobba/golang-siwes-training/week24/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ svc *service.InventoryService }

func NewHandler(svc *service.InventoryService) *Handler { return &Handler{svc: svc} }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrInsufficientStock):
		status = http.StatusConflict
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func parseID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%w: invalid id", domain.ErrInvalidInput)
	}
	return id, nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ps, err := h.svc.ListProducts(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ps)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	p, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name       string  `json:"name"`
		SKU        string  `json:"sku"`
		CategoryID int     `json:"category_id"`
		Price      float64 `json:"price"`
		Stock      int     `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, fmt.Errorf("%w: bad JSON body", domain.ErrInvalidInput))
		return
	}
	p := &domain.Product{Name: in.Name, SKU: in.SKU, CategoryID: in.CategoryID, Price: in.Price, Stock: in.Stock}
	if err := h.svc.CreateProduct(r.Context(), p); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var in struct {
		Delta int `json:"delta"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, fmt.Errorf("%w: bad JSON body", domain.ErrInvalidInput))
		return
	}
	p, err := h.svc.AdjustStock(r.Context(), id, in.Delta)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.svc.DeleteProduct(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
'@ | Set-Content -Encoding ascii internal/transport/http/handler.go

# 7. gRPC handler
@'
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
'@ | Set-Content -Encoding ascii internal/transport/grpc/handler.go

# 8. Server entry point (REST + gRPC + graceful shutdown)
@'
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
	"github.com/OluwaTobba/golang-siwes-training/week24/internal/repository"
	"github.com/OluwaTobba/golang-siwes-training/week24/internal/service"
	grpchandler "github.com/OluwaTobba/golang-siwes-training/week24/internal/transport/grpc"
	httphandler "github.com/OluwaTobba/golang-siwes-training/week24/internal/transport/http"
	pb "github.com/OluwaTobba/golang-siwes-training/week24/proto"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// Database (retry, because the DB container may start after the API)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env("DB_HOST", "localhost"), env("DB_PORT", "5432"),
		env("DB_USER", "postgres"), env("DB_PASSWORD", "secret"), env("DB_NAME", "siwes_db"))

	var db *gorm.DB
	var err error
	for i := 1; i <= 15; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("database not ready (%d/15): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("db: ", err)
	}
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		log.Fatal("migrate: ", err)
	}
	sqlDB, _ := db.DB()

	// Dependency injection
	repo := repository.NewGORMProductRepo(db)
	svc := service.NewInventoryService(repo)
	handler := httphandler.NewHandler(svc)

	// REST API
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer, middleware.RequestID)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/readyz", func(w http.ResponseWriter, req *http.Request) {
		if err := sqlDB.PingContext(req.Context()); err != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/api/v1/products", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.Get)
		r.Patch("/{id}/stock", handler.AdjustStock)
		r.Delete("/{id}", handler.Delete)
	})
	httpSrv := &http.Server{Addr: ":8080", Handler: r}

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("grpc listen: ", err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcSrv, grpchandler.NewGRPCHandler(svc))

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("REST :8080 | gRPC :50051 - ready")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
}
'@ | Set-Content -Encoding ascii cmd/server/main.go

# 9. gRPC test client
@'
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
'@ | Set-Content -Encoding ascii cmd/client/main.go

# 10. Dependencies
gofmt -w cmd internal proto
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
go get github.com/go-chi/chi/v5@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get github.com/prometheus/client_golang@latest
go mod tidy
go mod edit "-toolchain=none"
go vet ./...

# 11. Docker files (Go version in the image matches go.mod)
$gv = (Select-String -Path go.mod -Pattern '^go (\d+\.\d+)').Matches[0].Groups[1].Value
Write-Host "Using golang:$gv-alpine for the build image"

@'
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server
USER 1001:1001
EXPOSE 8080 50051
ENTRYPOINT ["/server"]
'@ | Set-Content -Encoding ascii Dockerfile

@'
.git
.gitignore
*.exe
*_test.go
.idea/
.vscode/
*.md
k8s/
'@ | Set-Content -Encoding ascii .dockerignore

@'
services:
  api:
    build:
      context: .
      args:
        GO_VERSION: ${GO_VERSION:-1.26}
    image: inventory-api:v1.0.0
    ports:
      - "8080:8080"
      - "50051:50051"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=secret
      - DB_NAME=siwes_db
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
    read_only: true
    tmpfs:
      - /tmp

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: siwes_db
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  pgdata:
'@ | Set-Content -Encoding ascii docker-compose.yml

# 12. Build and start
$env:GO_VERSION = $gv
docker compose up --build -d
Start-Sleep -Seconds 12
docker compose ps

# 13. Test REST
$base = "http://localhost:8080"
Invoke-RestMethod "$base/readyz" | Out-Null
$sku = "LAP-" + (Get-Date -Format HHmmss)
$new = Invoke-RestMethod -Method Post "$base/api/v1/products" -ContentType "application/json" -Body (@{name="Laptop"; sku=$sku; price=1500; stock=10} | ConvertTo-Json)
"Created: id=$($new.id) sku=$($new.sku) stock=$($new.stock)"
$adj = Invoke-RestMethod -Method Patch "$base/api/v1/products/$($new.id)/stock" -ContentType "application/json" -Body '{"delta":-3}'
"After -3: stock=$($adj.stock)"
Invoke-RestMethod "$base/api/v1/products" | Format-Table id, name, sku, price, stock
try { Invoke-RestMethod -Method Patch "$base/api/v1/products/$($new.id)/stock" -ContentType "application/json" -Body '{"delta":-100}' }
catch { "Oversell rejected with HTTP " + [int]$_.Exception.Response.StatusCode }
try { Invoke-RestMethod "$base/api/v1/products/99999" }
catch { "Unknown product -> HTTP " + [int]$_.Exception.Response.StatusCode }
"Metrics sample:"
((Invoke-WebRequest -UseBasicParsing "$base/metrics").Content -split "`n" | Select-Object -First 4)

# 14. Test gRPC
go run ./cmd/client

# 15. Server log
docker compose logs api --tail 8