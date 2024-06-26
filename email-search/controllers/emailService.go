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

func GetEmails(text string) (models.ZincResponse, error) {

	requestBody := models.SearchQueryRequest{
		SearchType: "match",
		Query: models.SearchQuery{
			Term:      text,
			Field:     "_all",
			StartTime: "2000-06-02T14:28:31.894Z",
			EndTime:   "2030-12-02T15:28:31.894Z",
		},
		SortFields: []string{"date"},
		From:       0,
		MaxResults: 2000000,
		Highlight: map[string]interface{}{
			"pre_tags":  []string{"<strong>"},
			"post_tags": []string{"</strong>"},
			"fields": map[string]interface{}{
				"content": map[string]interface{}{
					"pre_tags":  []string{},
					"post_tags": []string{},
				},
			},
		},
	}

	var zincResponse models.ZincResponse

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
		return zincResponse, err
	}

	apiEndpoint := os.Getenv("ZINC_SEARCH_SERVER_URL") + "api/" + os.Getenv("INDEX_NAME") + "/_search"

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return zincResponse, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("ZINC_FIRST_ADMIN_USER"), os.Getenv("ZINC_FIRST_ADMIN_PASSWORD"))

	res, err := http.DefaultClient.Do(req)
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
			Id:        hit.ID,
			From:      hit.Source.From,
			To:        hit.Source.To,
			Date:      hit.Source.Date,
			Subject:   hit.Source.Subject,
			Content:   hit.Source.Content,
			Highlight: hit.Highlight.Content,
		}

		// if len(hit.SortFields) > 0 {
		// 	fmt.Println("SortFields:", hit.SortFields)
		// }

		emails = append(emails, email)
	}

	return emails
}
