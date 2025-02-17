package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type RequestBody struct {
	RequireExport bool       `json:"requireExport"`
	ShopIds       []string   `json:"shopId"`
	RoundDate     string     `json:"roundDate"`
	Pagination    Pagination `json:"pagination"`
}

func main() {
	// Open the CSV file
	file, err := os.Open("lotto_shop.csv")
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip the header row
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	// Read all rows from the CSV file and get the first 3 shop_ids
	var shopIds []string
	for i := 0; i < 600; i++ { // Limiting to the first 3 entries
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV file: %v", err)
		}
		// Extract the shop_id from each row and add it to the slice
		shopIds = append(shopIds, record[1]) // Assuming shop_id is in the second column
	}

	// Construct the API payload with the first 3 shopIds
	payload := RequestBody{
		RequireExport: false,
		ShopIds:       shopIds,
		RoundDate:     "2024-09-01", // Adjust as needed
		Pagination: Pagination{
			Page:    1,
			PerPage: 10,
		},
	}

	// Convert payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Record the start time
	startTime := time.Now()

	// Send the POST request with the updated URL
	url := "http://core-lt-quota-manage.core-lt.svc.cluster.local:8080/core/quota/api/v1/search-info-shop"
	client := &http.Client{
		Timeout: time.Second * 120, // เพิ่มเวลา timeout เป็น 2 นาที
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Record the end time and calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Log the response time
	log.Printf("Total response time: %v", elapsedTime)

	// Read the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("API request successful")
	} else {
		fmt.Printf("API request failed with status: %s\n", resp.Status)
	}
}
