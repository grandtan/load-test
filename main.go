package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil" // เพิ่มการใช้งาน ioutil
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

	// Read all rows from the CSV file and get the first 250000 shop_ids
	var shopIds []string
	for i := 0; i < 250000; i++ {
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

	// Construct the API payload with the shopIds
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
	url := "http://core-lt-quota-manage.loadtest.svc:8080/core/quota/api/v1/search-info-shop"
	client := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
			MaxIdleConns:        10,               // Max number of idle connections
			MaxIdleConnsPerHost: 10,               // Max idle connections per host
			DisableKeepAlives:   false,            // Disable keep-alive connections
		},
		Timeout: time.Second * 120,
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

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response status code and body
	log.Printf("Response Status Code: %d", resp.StatusCode)
	log.Printf("Response Body: %s", string(body))

	// Check the response status
	if resp.StatusCode == http.StatusOK {
		fmt.Println("API request successful")
	} else {
		fmt.Printf("API request failed with status: %s\n", resp.Status)
	}
}
