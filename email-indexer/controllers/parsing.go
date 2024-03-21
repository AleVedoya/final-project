package controllers

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"trucode/finalproject/models"
)

func GetEmailsDir() ([]models.Email, error) {
	log.Println("Getting emails directory")

	emailsDir := "../data/enron_mail_20110402/maildir/bailey-s"
	// emailsDir := "../data/enron_mail_20110402/maildir/"
	var wg sync.WaitGroup
	var mu sync.Mutex
	var records []models.Email

	err := filepath.Walk(emailsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking path:", err)
			return err
		}
		if !info.IsDir() {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				emailData, err := processFile(file)
				if err != nil {
					log.Println("Error processing file:", err)
					return
				}
				mu.Lock()
				records = append(records, emailData)
				mu.Unlock()
			}(path)
		}
		return nil
	})

	if err != nil {
		log.Println("Error walking through directory:", err)
		return nil, err
	}

	wg.Wait()

	if err := CreateIndex(records); err != nil {
		log.Fatalf("failed to create index: %v", err)
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
