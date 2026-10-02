package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var (
	notes  = make(map[int]Note)
	nextID = 1
	mu     sync.RWMutex
)

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s  %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Handlers
func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.RLock()
		list := make([]Note, 0, len(notes))
		for _, n := range notes { list = append(list, n) }
		mu.RUnlock()
		jsonResponse(w, http.StatusOK, list)

	case http.MethodPost:
		var n Note
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		mu.Lock()
		n.ID = nextID
		nextID++
		notes[n.ID] = n
		mu.Unlock()
		jsonResponse(w, http.StatusCreated, n)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func noteByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		mu.RLock()
		n, ok := notes[id]
		mu.RUnlock()
		if !ok {
			jsonResponse(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		jsonResponse(w, http.StatusOK, n)

	case http.MethodDelete:
		mu.Lock()
		delete(notes, id)
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/notes", notesHandler)
	mux.HandleFunc("/notes/", noteByIDHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handler := loggingMiddleware(mux)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}