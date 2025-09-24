package main

import (
	"fmt"
	"net/http/httptest"

	"github.com/Zhukek/metrics/internal/handler"
	"github.com/Zhukek/metrics/internal/repository/inmemory"
	"github.com/go-resty/resty/v2"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func main() {
	// Create storage and server
	storage, _ := inmemory.NewStorage("", 0, false)
	server := httptest.NewServer(handler.NewRouter(storage, nil))
	defer server.Close()

	// Test with explicit gzip support
	fmt.Println("=== Testing GZIP compression ===")

	// First add some data
	httpc := resty.New().SetHostURL(server.URL)
	
	// Add a counter
	counter := int64(12345)
	_, err := httpc.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(&Metrics{
			ID:    "test_counter",
			MType: "counter",
			Delta: &counter,
		}).
		Post("update/")
	
	if err != nil {
		fmt.Printf("Error adding counter: %v\n", err)
		return
	}

	// Now test getting with gzip
	var result Metrics
	resp, err := httpc.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(&Metrics{
			ID:    "test_counter",
			MType: "counter",
		}).
		SetResult(&result).
		Post("value/")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode())
	fmt.Printf("Content-Type: %s\n", resp.Header().Get("Content-Type"))
	fmt.Printf("Content-Encoding: %s\n", resp.Header().Get("Content-Encoding"))
	fmt.Printf("Accept-Encoding sent: gzip\n")
	fmt.Printf("Response size: %d bytes\n", len(resp.Body()))
	
	if resp.Header().Get("Content-Encoding") == "gzip" {
		fmt.Println("✅ GZIP compression working!")
	} else {
		fmt.Println("❌ GZIP compression not working")
		fmt.Println("Response body:", string(resp.Body()))
	}

	// Test with larger response (list all metrics)
	fmt.Println("\n=== Testing GZIP with larger response ===")
	resp2, err := httpc.R().
		SetHeader("Accept-Encoding", "gzip").
		Get("/")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("List status: %d\n", resp2.StatusCode())
	fmt.Printf("List Content-Type: %s\n", resp2.Header().Get("Content-Type"))
	fmt.Printf("List Content-Encoding: %s\n", resp2.Header().Get("Content-Encoding"))
	fmt.Printf("List response size: %d bytes\n", len(resp2.Body()))
}
