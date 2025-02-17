package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

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

	// Create a slice to store shop IDs
	var shopIds []string

	// Read all rows from the CSV file
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		// Assuming the shop_id is in the second column (index 1)
		if len(record) > 1 {
			shopIds = append(shopIds, record[1]) // Add shop_id to the slice
		}
	}

	// Prepare the curl command format
	shopIdsStr := strings.Join(shopIds, `","`) // Join shopIds as a comma-separated list
	curlCommand := fmt.Sprintf(`
curl --location 'http://core-lt-quota-manage.core-lt.svc.cluster.local:8080/core/quota/api/v1/search-info-shop' \
--header 'Content-Type: application/json' \
--data '{
        "requireExport": false,
        "shopId": ["%s"],
        "roundDate": "2024-09-01",
        "pagination": {
          "page": 1,
          "perPage": 10
        }
      }'`, shopIdsStr)

	// Write the curl command to a file
	outputFile := "generated_curl_command.txt"
	err = os.WriteFile(outputFile, []byte(curlCommand), 0644)
	if err != nil {
		log.Fatalf("Error writing curl command to file: %v", err)
	}

	// Confirm the file creation
	fmt.Printf("Curl command successfully generated: %s\n", outputFile)
}
