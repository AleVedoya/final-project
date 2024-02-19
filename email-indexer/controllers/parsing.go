// controllers/parsing.go
package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"trucode/finalproject/models"
)

func GetEmailsDir() ([]models.Email, error) {
	log.Println("Getting emalis directory")

	emailsDir := "../data/enron_mail_20110402/maildir/"
	var records []models.Email

	log.Printf("Emails Dir: %s", emailsDir)

	err := filepath.Walk(emailsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking path:", err)
			return err
		}

		if !info.IsDir() {
			emailData, err := processFile(path)
			if err != nil {
				log.Println(err)
				return nil
			}
			records = append(records, emailData)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking through directory: %v", err)
	}

	return records, nil
}

func processFile(path string) (models.Email, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return models.Email{}, err
	}

	email := parseEmail(string(content))
	return email, nil
}

func parseEmail(rawEmail string) models.Email {
	var email models.Email

	headersAndBody := strings.SplitN(rawEmail, "\n\n", 2)
	if len(headersAndBody) < 2 {
		return email
	}

	headers := headersAndBody[0]
	body := headersAndBody[1]

	headerLines := strings.Split(headers, "\n")
	header := make(map[string]string)
	for _, line := range headerLines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			header[parts[0]] = parts[1]
		}
	}

	recipients := strings.Split(header["To"], ", ")
	email = models.Email{
		Message_ID: header["Message-ID"],
		Date:       header["Date"],
		From:       header["From"],
		To:         recipients,
		Subject:    header["Subject"],
		Content:    body,
	}

	return email
}
