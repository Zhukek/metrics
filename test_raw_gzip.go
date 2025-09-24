package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Zhukek/metrics/internal/handler"
	"github.com/Zhukek/metrics/internal/repository/inmemory"
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

	// Add a counter first
	counter := int64(12345)
	metric := &Metrics{
		ID:    "test_counter",
		MType: "counter",
		Delta: &counter,
	}
	
	body, _ := json.Marshal(metric)
	
	// Add the counter
	req, _ := http.NewRequest("POST", server.URL+"/update/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error adding counter: %v\n", err)
		return
	}
	resp.Body.Close()
	
	fmt.Printf("Add counter status: %d\n", resp.StatusCode)

	// Now test getting with raw http client and gzip
	getMetric := &Metrics{
		ID:    "test_counter",
		MType: "counter",
	}
	
	getBody, _ := json.Marshal(getMetric)
	
	req, _ = http.NewRequest("POST", server.URL+"/value/", bytes.NewBuffer(getBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")  // Request gzip compression
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("\n=== Raw HTTP Client Test ===\n")
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("Content-Encoding: %s\n", resp.Header.Get("Content-Encoding"))
	fmt.Printf("Content-Length: %s\n", resp.Header.Get("Content-Length"))
	
	if resp.Header.Get("Content-Encoding") == "gzip" {
		fmt.Println("✅ GZIP compression working with raw HTTP client!")
	} else {
		fmt.Println("❌ GZIP compression not working with raw HTTP client")
	}

	// Read the raw response body (will be gzipped if compression is working)
	bodyBytes := make([]byte, 1024)
	n, _ := resp.Body.Read(bodyBytes)
	bodyBytes = bodyBytes[:n]
	
	fmt.Printf("Raw response body length: %d bytes\n", len(bodyBytes))
	fmt.Printf("First few bytes (hex): %x\n", bodyBytes[:min(10, len(bodyBytes))])
	
	// Check if it looks like gzip (starts with 1f 8b)
	if len(bodyBytes) >= 2 && bodyBytes[0] == 0x1f && bodyBytes[1] == 0x8b {
		fmt.Println("✅ Response body is gzip compressed!")
	} else {
		fmt.Println("❌ Response body is not gzip compressed")
		fmt.Printf("Response body as string: %s\n", string(bodyBytes))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
