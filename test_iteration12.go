package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

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
	// Create storage and server (without hasher for this test)
	storage, _ := inmemory.NewStorage("", 0, false)
	server := httptest.NewServer(handler.NewRouter(storage, nil))
	defer server.Close()

	// Create HTTP client with gzip support
	httpc := resty.New().SetHostURL(server.URL).
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	
	idCounter := "CounterBatchZip" + strconv.Itoa(rnd.Intn(256))
	idGauge := "GaugeBatchZip" + strconv.Itoa(rnd.Intn(256))
	valueCounter1, valueCounter2 := int64(rnd.Int31()), int64(rnd.Int31())
	var valueCounter0 int64
	valueGauge1, valueGauge2 := rnd.Float64()*1e6, rnd.Float64()*1e6

	fmt.Printf("Testing with Counter ID: %s, Gauge ID: %s\n", idCounter, idGauge)
	fmt.Printf("Counter values: %d, %d\n", valueCounter1, valueCounter2)
	fmt.Printf("Gauge values: %f, %f\n", valueGauge1, valueGauge2)

	// Test 1: Get random counter (should be 0 or not found)
	fmt.Println("\n=== Test 1: Get initial counter value ===")
	var result Metrics
	req := httpc.R()
	resp, err := req.
		SetBody(&Metrics{
			ID:    idCounter,
			MType: "counter",
		}).
		SetResult(&result).
		Post("value/")

	if err != nil {
		fmt.Printf("❌ Error getting counter: %v\n", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode())
		fmt.Printf("Content-Type: %s\n", resp.Header().Get("Content-Type"))
		fmt.Printf("Content-Encoding: %s\n", resp.Header().Get("Content-Encoding"))
		
		switch resp.StatusCode() {
		case http.StatusOK:
			if result.Delta != nil {
				valueCounter0 = *result.Delta
				fmt.Printf("✅ Initial counter value: %d\n", valueCounter0)
			}
		case http.StatusNotFound:
			fmt.Printf("✅ Counter not found (expected for new counter)\n")
			valueCounter0 = 0
		default:
			fmt.Printf("❌ Unexpected status: %d\n", resp.StatusCode())
		}
	}

	// Test 2: Batch update metrics
	fmt.Println("\n=== Test 2: Batch update metrics ===")
	metrics := []Metrics{
		{
			ID:    idCounter,
			MType: "counter",
			Delta: &valueCounter1,
		},
		{
			ID:    idGauge,
			MType: "gauge",
			Value: &valueGauge1,
		},
		{
			ID:    idCounter,
			MType: "counter",
			Delta: &valueCounter2,
		},
		{
			ID:    idGauge,
			MType: "gauge",
			Value: &valueGauge2,
		},
	}

	req = httpc.R().
		SetHeader("Content-Type", "application/json")
	resp, err = req.SetBody(metrics).
		Post("updates/")

	if err != nil {
		fmt.Printf("❌ Error in batch update: %v\n", err)
	} else {
		fmt.Printf("Batch update status: %d\n", resp.StatusCode())
		fmt.Printf("Content-Type: %s\n", resp.Header().Get("Content-Type"))
		if resp.StatusCode() == http.StatusOK {
			fmt.Printf("✅ Batch update successful\n")
		} else {
			fmt.Printf("❌ Batch update failed with status: %d\n", resp.StatusCode())
		}
	}

	// Test 3: Check counter value
	fmt.Println("\n=== Test 3: Check counter value ===")
	var counterResult Metrics
	req = httpc.R()
	resp, err = req.
		SetBody(&Metrics{
			ID:    idCounter,
			MType: "counter",
		}).
		SetResult(&counterResult).
		Post("value/")

	if err != nil {
		fmt.Printf("❌ Error getting counter: %v\n", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode())
		fmt.Printf("Content-Type: %s\n", resp.Header().Get("Content-Type"))
		fmt.Printf("Content-Encoding: %s\n", resp.Header().Get("Content-Encoding"))
		
		if resp.StatusCode() == http.StatusOK {
			if counterResult.Delta != nil {
				expectedCounter := valueCounter0 + valueCounter1 + valueCounter2
				if *counterResult.Delta == expectedCounter {
					fmt.Printf("✅ Counter value correct: %d (expected: %d)\n", *counterResult.Delta, expectedCounter)
				} else {
					fmt.Printf("❌ Counter value mismatch: got %d, expected %d\n", *counterResult.Delta, expectedCounter)
				}
			} else {
				fmt.Printf("❌ Counter Delta is nil\n")
			}
		} else {
			fmt.Printf("❌ Failed to get counter value, status: %d\n", resp.StatusCode())
		}
	}

	// Test 4: Check gauge value
	fmt.Println("\n=== Test 4: Check gauge value ===")
	var gaugeResult Metrics
	req = httpc.R()
	resp, err = req.
		SetBody(&Metrics{
			ID:    idGauge,
			MType: "gauge",
		}).
		SetResult(&gaugeResult).
		Post("value/")

	if err != nil {
		fmt.Printf("❌ Error getting gauge: %v\n", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode())
		fmt.Printf("Content-Type: %s\n", resp.Header().Get("Content-Type"))
		fmt.Printf("Content-Encoding: %s\n", resp.Header().Get("Content-Encoding"))
		
		if resp.StatusCode() == http.StatusOK {
			if gaugeResult.Value != nil {
				if *gaugeResult.Value == valueGauge2 {
					fmt.Printf("✅ Gauge value correct: %f (expected: %f)\n", *gaugeResult.Value, valueGauge2)
				} else {
					fmt.Printf("❌ Gauge value mismatch: got %f, expected %f\n", *gaugeResult.Value, valueGauge2)
				}
			} else {
				fmt.Printf("❌ Gauge Value is nil\n")
			}
		} else {
			fmt.Printf("❌ Failed to get gauge value, status: %d\n", resp.StatusCode())
		}
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ All tests completed. Check results above.")
}
