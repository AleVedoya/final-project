package models

type Email struct {
	Id        string   `json:"id"`
	From      string   `json:"from"`
	To        []string   `json:"to"`
	Content   string   `json:"content"`
	Subject   string   `json:"subject"`
	Date      string   `json:"date"`
}

type ZincResponse struct {
	Took   int     `json:"took"`
	Emails []Email `json:"emails"`
}

type EmailSearchResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Index     string  `json:"_index"`
			Type      string  `json:"_type"`
			ID        string  `json:"_id"`
			Score     float64 `json:"_score"`
			Timestamp string  `json:"@timestamp"`
			Source    Email   `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
