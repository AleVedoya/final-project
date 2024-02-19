package models

type Email struct {
	Message_ID string   `json:"message_id"`
	Date       string   `json:"date"`
	From       string   `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	Content    string   `json:"content"`
	Filepath   string   `json:"filepath"`
}
