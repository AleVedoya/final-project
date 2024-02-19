package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"trucode/finalproject/models"
)

func CheckIfIndexExists() *http.Response {
	// http://localhost:4080/api/enron/_mapping
	log.Printf("Checking if index %s exists", os.Getenv("INDEX_NAME"))
	url := os.Getenv("ZINC_SEARCH_SERVER_URL") + "api/" + os.Getenv("INDEX_NAME") + "/_mapping"

	req, err := makeRequestWithAuth("GET", url, "")
	if err != nil {
		log.Fatal(err)
	}

	zincUser := os.Getenv("ZINC_FIRST_ADMIN_USER")
	zincPass := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")
	req.SetBasicAuth(zincUser, zincPass)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == 400 {
		log.Printf("Index %s does not exist", os.Getenv("INDEX_NAME"))
		return res
	}

	if res.StatusCode == 200 {
		log.Printf("Index %s does exist.", os.Getenv("INDEX_NAME"))
	}

	log.Printf("Response status code: %d", res.StatusCode)
	log.Print("Index checks finalized")

	return res
}

func CreateIndex(records []models.Email) error {
	log.Printf("Creating %s index:", os.Getenv("INDEX_NAME"))

	// http://localhost:4080/api/_bulkv2
	bulkURL := fmt.Sprintf(os.Getenv("ZINC_SEARCH_SERVER_URL") + "api/_bulkv2")
	indexName := os.Getenv("INDEX_NAME")

	bulk := struct {
		IndexName string         `json:"index"`
		Records   []models.Email `json:"records"`
	}{
		IndexName: indexName,
		Records:   records,
	}

	jsonBulk, err := json.Marshal(bulk)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := makeRequestWithAuth("POST", bulkURL, string(jsonBulk))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()

	log.Printf("HTTP response status code: %d", res.StatusCode)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create index, status code: %d", res.StatusCode)
	}

	log.Println("Index created successfully")

	return nil
}

func DeleteIndex(indexName string) error {
	// http://localhost:4080/api/index/enron
	url := os.Getenv("ZINC_SEARCH_SERVER_URL") + "api/index/" + os.Getenv("INDEX_NAME")

	req, err := makeRequestWithAuth("DELETE", url, "")
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete indexer, status code: %d", res.StatusCode)
	}

	log.Println("Index deleted successfully")

	return nil
}

func makeRequestWithAuth(method string, url string, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	zincUser := os.Getenv("ZINC_FIRST_ADMIN_USER")
	zincPass := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")
	req.SetBasicAuth(zincUser, zincPass)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
