package models

type Email struct {
	Id        string   `json:"id"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Date      string   `json:"date"`
	Subject   string   `json:"subject"`
	Content   string   `json:"content"`
	Highlight []string `json:"highlight"`
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
			Index      string   `json:"_index"`
			Type       string   `json:"_type"`
			ID         string   `json:"_id"`
			Score      float64  `json:"_score"`
			Timestamp  string   `json:"@timestamp"`
			Source     Email    `json:"_source"`
			SortFields []string `json:"sortFields"`
			Highlight  struct {
				Content []string `json:"content"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

type SearchQueryRequest struct {
	SearchType string                 `json:"search_type"`
	Query      SearchQuery            `json:"query"`
	From       int                    `json:"from"`
	MaxResults int                    `json:"max_results"`
	SortFields []string               `json:"sort_fields"`
	Source     map[string]interface{} `json:"_source"`
	Highlight  map[string]interface{} `json:"highlight"`
}

type SearchQuery struct {
	Term      string `json:"term"`
	Field     string `json:"field"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
