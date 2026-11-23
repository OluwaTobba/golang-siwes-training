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