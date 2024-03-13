package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Defines receipt structure.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Defines item structure in a receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// For points response.
type PointsResponse struct {
	Points int `json:"points"`
}

// For ID response after processing receipt.
type IDResponse struct {
	ID string `json:"id"`
}

var (
	mu          sync.Mutex // Protects receiptsMap.
	receiptsMap = make(map[string]int) // Stores receipts and points.
)

func main() {
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported method.", http.StatusNotFound)
		return
	}

	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String() // Generate unique ID for receipt.

	mu.Lock()
	receiptsMap[id] = calculatePoints(receipt) // Calculate and store points.
	mu.Unlock()

	json.NewEncoder(w).Encode(IDResponse{ID: id})
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points") // Extract ID.

	mu.Lock()
	points, exists := receiptsMap[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Receipt not found.", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(PointsResponse{Points: points})
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// Alphanumeric characters in retailer name.
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	points += len(reg.ReplaceAllString(receipt.Retailer, ""))

	// Round dollar bonus.
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Multiple of 0.25 bonus.
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Points per item pair.
	points += (len(receipt.Items) / 2) * 5

	// Bonus for item descriptions and odd day.
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// Odd day bonus.
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if date.Day()%2 != 0 {
			points += 6
		}
	}

	// Time bonus.
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
			points += 10
		}
	}

	return points
}
