package main

import (
	"context"
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

	"github.com/yourname/siwes/week24/internal/domain"
	"github.com/yourname/siwes/week24/internal/repository"
	"github.com/yourname/siwes/week24/internal/service"
	httphandler "github.com/yourname/siwes/week24/internal/transport/http"
	grpchandler "github.com/yourname/siwes/week24/internal/transport/grpc"
	pb "github.com/yourname/siwes/week24/proto"
)

func main() {
	// ── Database ───────────────────────────────────────────────
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { log.Fatal("db:", err) }
	db.AutoMigrate(&domain.Product{})

	// ── Dependency injection ───────────────────────────────────
	repo    := repository.NewGORMProductRepo(db)
	svc     := service.NewInventoryService(repo)
	handler := httphandler.NewHandler(svc)

	// ── REST API ───────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer, middleware.RequestID)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	r.Get("/readyz",  func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/api/v1/products", func(r chi.Router) {
		r.Get("/",        handler.List)
		r.Post("/",       handler.Create)
		r.Get("/{id}",    handler.Get)
		r.Patch("/{id}/stock", handler.AdjustStock)
		r.Delete("/{id}", handler.Delete)
	})
	httpSrv := &http.Server{Addr: ":8080", Handler: r}

	// ── gRPC Server ───────────────────────────────────────────
	lis, _ := net.Listen("tcp", ":50051")
	grpcSrv := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcSrv, grpchandler.NewGRPCHandler(svc))

	// ── Graceful shutdown ─────────────────────────────────────
	go func() { log.Fatal(httpSrv.ListenAndServe()) }()
	go func() { log.Fatal(grpcSrv.Serve(lis)) }()
	log.Println("REST :8080  |  gRPC :50051  —  ready")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
}