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