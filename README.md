# 🐹 SIWES Golang Training — 24-Week Industrial Attachment

> **Student Industrial Work Experience Scheme (SIWES)**  
> Department: **DevOps & Software Engineering**  
> Duration: **24 Weeks**  
> Technology: **Go (Golang) 1.22**

---

## 📋 Overview

This repository contains all code, projects, and notes produced during a 24-week SIWES industrial training attachment focused on the **Go programming language**. The training followed a structured weekly curriculum — progressing from language fundamentals through to production-grade microservice development, containerisation, and Kubernetes deployment.

Each week has its own folder (`week01/` → `week24/`) containing working, commented Go source files (and supporting config files) that correspond to that week's learning objectives.

---

## 🗂️ Repository Structure

```
siwes-golang-training/
├── week01/          # Environment setup & Hello World
├── week02/          # Variables, data types & constants
├── week03/          # Control flow — for, if/else, switch
├── week04/          # Functions, multiple returns & defer
├── week05/          # Arrays, slices & maps
├── week06/          # Structs, methods & embedding
├── week07/          # Interfaces & polymorphism
├── week08/          # Error handling — errors, panic & recover
├── week09/          # Goroutines & channels
├── week10/          # Advanced concurrency — select, sync & context
├── week11/          # Packages, modules & standard library
├── week12/          # File I/O, JSON & CSV
├── week13/          # HTTP server with net/http
├── week14/          # REST API with chi router
├── week15/          # JSON, XML & data serialisation
├── week16/          # database/sql with PostgreSQL
├── week17/          # GORM ORM
├── week18/          # Testing — unit, table-driven & benchmarks
├── week19/          # JWT authentication & middleware
├── week20/          # Docker — multi-stage builds & Compose
├── week21/          # CLI tools with Cobra & Viper
├── week22/          # gRPC & Protocol Buffers
├── week23/          # Kubernetes deployment
├── week24/          # Capstone — Inventory Management Microservice
└── README.md
```

---

## 📅 Weekly Curriculum

| Week | Topic | Key Concepts |
|------|-------|-------------|
| 01 | Go Environment Setup | Toolchain, `go run`, `go build`, `go mod init` |
| 02 | Variables & Data Types | Primitive types, `:=`, zero values, `iota` |
| 03 | Control Flow | `for`, `range`, `if/else`, `switch` |
| 04 | Functions | Multiple returns, named returns, `defer`, variadic |
| 05 | Collections | Arrays, slices, maps, 2D slices |
| 06 | Structs & Methods | Embedding, value/pointer receivers, JSON tags |
| 07 | Interfaces | Implicit satisfaction, polymorphism, type switch |
| 08 | Error Handling | Custom errors, `errors.As/Is`, wrapping, `panic/recover` |
| 09 | Goroutines & Channels | Pipelines, fan-out, buffered channels |
| 10 | Advanced Concurrency | `select`, `sync.WaitGroup`, `Mutex`, `context` |
| 11 | Packages & Modules | `go.mod`, standard library, custom packages |
| 12 | File I/O | `os`, `bufio`, `encoding/json`, `encoding/csv` |
| 13 | HTTP Server | `net/http`, middleware, in-memory CRUD |
| 14 | REST API | `chi` router, layered architecture, validation |
| 15 | Serialisation | Custom JSON marshaller, XML, CSV transforms |
| 16 | Database (Raw SQL) | `database/sql`, `pgx`, transactions, connection pool |
| 17 | GORM ORM | Models, associations, hooks, `AutoMigrate` |
| 18 | Testing | Table-driven tests, `httptest`, benchmarks |
| 19 | Authentication | JWT, auth middleware, RBAC |
| 20 | Docker | Multi-stage builds, Compose, `.dockerignore` |
| 21 | CLI Tools | Cobra, Viper, nested subcommands, flags |
| 22 | gRPC | Protocol Buffers, streaming RPC, interceptors |
| 23 | Kubernetes | Deployment, Service, ConfigMap, HPA, probes |
| 24 | Capstone Project | Full Inventory Management Microservice |

---

## 🚀 Capstone Project (Week 24)

A production-grade **Inventory Management Microservice** demonstrating all skills from the 24-week programme:

### Features
- ✅ Clean architecture — `cmd/`, `internal/`, `pkg/`, `api/` layers
- ✅ REST API — `chi` router, JWT authentication, RBAC middleware
- ✅ gRPC interface — for internal service-to-service queries
- ✅ PostgreSQL — GORM ORM with migrations and business rule enforcement
- ✅ Docker — multi-stage build (~10 MB final image)
- ✅ Kubernetes — Deployment, Service, HPA, health probes via Kustomize
- ✅ Prometheus metrics — `/metrics` endpoint via `promhttp`
- ✅ OpenAPI docs — generated with `swaggo/swag`
- ✅ Tests — unit + integration coverage

### Quick Start

```bash
# Clone the repo
git clone https://github.com/<your-username>/siwes-golang-training.git
cd siwes-golang-training/week24

# Run with Docker Compose (API + PostgreSQL)
docker compose up --build

# API will be available at http://localhost:8080
# gRPC will be available at localhost:50051
```

### API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/auth/login` | Get JWT token | Public |
| `GET` | `/api/v1/products` | List all products | Required |
| `POST` | `/api/v1/products` | Create product | Required |
| `GET` | `/api/v1/products/{id}` | Get product by ID | Required |
| `PATCH` | `/api/v1/products/{id}/stock` | Adjust stock | Required |
| `DELETE` | `/api/v1/products/{id}` | Delete product | Admin |
| `GET` | `/healthz` | Liveness probe | Public |
| `GET` | `/readyz` | Readiness probe | Public |
| `GET` | `/metrics` | Prometheus metrics | Public |

---

## 🛠️ Prerequisites

Make sure you have the following installed to run the examples:

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.22+ | Running all Go code |
| [Docker](https://docs.docker.com/get-docker/) | 24+ | Week 20 & capstone |
| [Docker Compose](https://docs.docker.com/compose/) | v2 | Multi-service stack |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | 1.29+ | Week 23 Kubernetes |
| [minikube](https://minikube.sigs.k8s.io/) | 1.32+ | Local K8s cluster |
| [protoc](https://grpc.io/docs/protoc-installation/) | 3.x | Week 22 gRPC codegen |
| PostgreSQL | 16 | Weeks 16–17, 24 |

---

## ▶️ Running Individual Weeks

Each week folder is a self-contained Go program or module.

```bash
# Example: Run week 3 (control flow)
cd week03
go run control.go

# Example: Run week 9 (goroutines)
cd week09
go run concurrency.go

# Example: Run the week 14 REST API
cd week14
go mod tidy
go run main.go
# → API running on :8080

# Example: Run tests for week 18
cd week18
go test ./... -v
go test -bench=. -benchmem
```

---

## 📦 Key Dependencies

| Package | Purpose | Used in |
|---------|---------|---------|
| [`github.com/go-chi/chi/v5`](https://github.com/go-chi/chi) | HTTP router | Weeks 14, 19, 24 |
| [`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx) | PostgreSQL driver | Week 16 |
| [`gorm.io/gorm`](https://gorm.io) | ORM | Weeks 17, 24 |
| [`github.com/golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt) | JWT auth | Week 19 |
| [`github.com/spf13/cobra`](https://github.com/spf13/cobra) | CLI framework | Week 21 |
| [`github.com/spf13/viper`](https://github.com/spf13/viper) | Config management | Week 21 |
| [`google.golang.org/grpc`](https://google.golang.org/grpc) | gRPC framework | Weeks 22, 24 |
| [`google.golang.org/protobuf`](https://google.golang.org/protobuf) | Protocol Buffers | Weeks 22, 24 |
| [`github.com/prometheus/client_golang`](https://github.com/prometheus/client_golang) | Metrics | Week 24 |

---

## 📚 Resources Used During Training

- [The Go Programming Language](https://www.gopl.io/) — Donovan & Kernighan
- [Go Official Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/) — Katherine Cox-Buday
- [100 Go Mistakes and How to Avoid Them](https://100go.co/) — Teiva Harsanyi
- [gRPC Official Docs](https://grpc.io/docs/languages/go/)
- [GORM Documentation](https://gorm.io/docs/)

---

## 👤 About

**Michael Ogundipe**  
SIWES Industrial Trainee  
Department of Computer Science  
National Open University of Nigeria 

---

## 📄 Licence

This repository is for educational purposes as part of the SIWES industrial training programme.  
All code is original work produced during the attachment period.

---

> *"The Go programming language is an open source project to make programmers more productive."*  
> — golang.org
