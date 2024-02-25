package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"trucode/search/models"
)

var (
	apiEndpoint = os.Getenv("ZINC_SEARCH_SERVER_URL") + "api/" + os.Getenv("INDEX_NAME") + "/_search"

	httpClient = &http.Client{}
)

func GetEmails(email string) (models.ZincResponse, error) {
	requestBody := map[string]interface{}{
		"search_type": "match",
		"query": map[string]interface{}{
			"term":  email,
			"field": "content",
		},
		"from":        0,
		"max_results": 200,
		"_source":     []string{},
	}

	var zincResponse models.ZincResponse

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
		return zincResponse, err
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return zincResponse, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("ZINC_FIRST_ADMIN_USER"), os.Getenv("ZINC_FIRST_ADMIN_PASSWORD"))

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return zincResponse, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return zincResponse, err
	}

	res.Body.Close()

	response := models.EmailSearchResult{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err)
		return zincResponse, err
	}

	zincResponse = models.ZincResponse{
		Took:   response.Hits.Total.Value,
		Emails: convertToEmails(response),
	}

	return zincResponse, nil
}

func convertToEmails(response models.EmailSearchResult) []models.Email {
	emails := []models.Email{}

	for _, hit := range response.Hits.Hits {
		email := models.Email{
			Id:      hit.ID,
			From:    hit.Source.From,
			To:      hit.Source.To,
			Subject: hit.Source.Subject,
			Content: hit.Source.Content,
			Date:    hit.Source.Date,
		}
		emails = append(emails, email)
	}

	return emails
}
