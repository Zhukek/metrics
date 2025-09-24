package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	compress "github.com/Zhukek/metrics/internal/middlewares"
)

func main() {
	// Create a simple handler that returns JSON
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test","type":"counter","delta":12345}`))
	})

	// Wrap with gzip middleware
	gzipHandler := compress.GzipMiddleware(handler)

	// Create test request with Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Call handler
	gzipHandler.ServeHTTP(w, req)
	
	// Check results
	fmt.Printf("Status: %d\n", w.Code)
	fmt.Printf("Content-Type: %s\n", w.Header().Get("Content-Type"))
	fmt.Printf("Content-Encoding: %s\n", w.Header().Get("Content-Encoding"))
	fmt.Printf("Accept-Encoding in request: %s\n", req.Header.Get("Accept-Encoding"))
	fmt.Printf("Response body length: %d\n", len(w.Body.Bytes()))
	fmt.Printf("Response body: %s\n", w.Body.String())
	
	// Check if Accept-Encoding contains gzip
	acceptEncoding := req.Header.Get("Accept-Encoding")
	supportsGzip := strings.Contains(acceptEncoding, "gzip")
	fmt.Printf("Supports gzip: %t\n", supportsGzip)
	
	if w.Header().Get("Content-Encoding") == "gzip" {
		fmt.Println("✅ GZIP compression working in middleware!")
	} else {
		fmt.Println("❌ GZIP compression not working in middleware")
	}
}
