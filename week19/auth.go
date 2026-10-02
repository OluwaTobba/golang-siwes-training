package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key-change-in-prod")

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type ctxKey string
const claimsKey ctxKey = "claims"

func GenerateToken(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil { return nil, err }
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid { return nil, fmt.Errorf("invalid token") }
	return claims, nil
}

// Auth middleware — validates JWT and injects claims into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized); return
		}
		claims, err := ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized); return
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RBAC middleware — restricts to given roles
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(claimsKey).(*Claims)
			if !ok { http.Error(w, "no claims", http.StatusForbidden); return }
			for _, role := range roles {
				if claims.Role == role { next.ServeHTTP(w, r); return }
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func main() {
	mux := http.NewServeMux()

	// Public — login
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		token, _ := GenerateToken(1, "admin")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	// Protected — any authenticated user
	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(claimsKey).(*Claims)
		json.NewEncoder(w).Encode(map[string]any{"user_id": claims.UserID, "role": claims.Role})
	})
	mux.Handle("/profile", AuthMiddleware(protected))

	// Admin only
	adminOnly := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome, Admin!")
	})
	mux.Handle("/admin", AuthMiddleware(RequireRole("admin")(adminOnly)))

	fmt.Println("Auth server on :8080")
	http.ListenAndServe(":8080", mux)
}