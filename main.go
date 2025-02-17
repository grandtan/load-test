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

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
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

	// Read all rows from the CSV file
	var shopIds []string
	for {
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

	// Chunk the shopIds into smaller pieces
	chunks := chunkSlice(shopIds, 1000) // Adjust chunk size as needed

	// Record the start time
	startTime := time.Now()

	// Send each chunk of shopIds as a separate request
	url := "http://core-lt-quota-manage.core-lt.svc.cluster.local:8080/core/quota/api/v1/search-info-shop"
	client := &http.Client{}
	for _, chunk := range chunks {
		// Construct the API payload for each chunk
		payload := RequestBody{
			RequireExport: false,
			ShopIds:       chunk,
			RoundDate:     "2024-09-01",
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

		// Send the POST request
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

		// Read the response
		if resp.StatusCode == http.StatusOK {
			fmt.Println("API request successful")
		} else {
			fmt.Printf("API request failed with status: %s\n", resp.Status)
		}
	}

	// Record the end time and calculate elapsed time
	elapsedTime := time.Since(startTime)
	log.Printf("Total response time: %v", elapsedTime)
}
