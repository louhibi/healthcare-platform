package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := NewLocationStore()
	h := NewHandler(store)
	r.GET("/api/locations/countries", h.GetCountries)
	r.GET("/api/locations/countries/:code/cities", h.GetCitiesByCountry)
	return r
}

func TestGetCountries(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/locations/countries", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" || body == "null" {
		t.Fatalf("expected non-empty countries list")
	}
}

func TestGetCitiesByCountry_ValidAndFilter(t *testing.T) {
	r := setupRouter()
	// Valid country
	req := httptest.NewRequest(http.MethodGet, "/api/locations/countries/US/cities", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Filter query 'yo' should match 'New York'
	req = httptest.NewRequest(http.MethodGet, "/api/locations/countries/US/cities?q=yo", nil)
	w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "New York") {
		t.Fatalf("expected filtered result to contain New York, got %s", w.Body.String())
	}
}

func TestGetCitiesByCountry_NotFound(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/locations/countries/ZZ/cities", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// naive contains to avoid extra deps
func contains(s, sub string) bool { return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && (indexOf(s, sub) >= 0))) }
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if s[i+j] != sub[j] { match = false; break }
		}
		if match { return i }
	}
	return -1
}
