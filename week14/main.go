package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ── Domain ────────────────────────────────────────────────────
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ── Repository (in-memory) ────────────────────────────────────
type ProductRepo struct {
	mu      sync.RWMutex
	items   map[int]Product
	counter int
}
func NewRepo() *ProductRepo { return &ProductRepo{items: make(map[int]Product)} }
func (r *ProductRepo) Create(p Product) Product {
	r.mu.Lock(); defer r.mu.Unlock()
	r.counter++; p.ID = r.counter; r.items[p.ID] = p; return p
}
func (r *ProductRepo) List() []Product {
	r.mu.RLock(); defer r.mu.RUnlock()
	list := make([]Product, 0, len(r.items))
	for _, p := range r.items { list = append(list, p) }; return list
}
func (r *ProductRepo) Get(id int) (Product, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	p, ok := r.items[id]
	if !ok { return Product{}, errors.New("not found") }
	return p, nil
}
func (r *ProductRepo) Delete(id int) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok { return errors.New("not found") }
	delete(r.items, id); return nil
}

// ── Handler ───────────────────────────────────────────────────
type Handler struct { repo *ProductRepo }

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, h.repo.List())
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); return
	}
	if p.Name == "" || p.Price <= 0 {
		respond(w, http.StatusBadRequest, map[string]string{"error": "name and price required"}); return
	}
	respond(w, http.StatusCreated, h.repo.Create(p))
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := h.repo.Get(id)
	if err != nil { respond(w, http.StatusNotFound, map[string]string{"error": err.Error()}); return }
	respond(w, http.StatusOK, p)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := h.repo.Delete(id); err != nil {
		respond(w, http.StatusNotFound, map[string]string{"error": err.Error()}); return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	h := &Handler{repo: NewRepo()}
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Get("/",        h.List)
			r.Post("/",       h.Create)
			r.Get("/{id}",    h.Get)
			r.Delete("/{id}", h.Delete)
		})
	})
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	fmt.Println("API running on :8080")
	http.ListenAndServe(":8080", r)
}