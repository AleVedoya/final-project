// controllers/parsing.go
package controllers

import (
	"bufio"
	"strings"
	"trucode/finalproject/models"
)

// ParseData is a placeholder function that parses data from a Scanner
// and returns an instance of IndexedData.
func ParseData(dataLines *bufio.Scanner, id int) models.Email {
	var data models.Email

	for dataLines.Scan() {
		line := dataLines.Text()

		data.ID = id

		if strings.Contains(line, "Message-ID:") {
			data.Message_ID = extractValue(line, "Message-ID:")
		} else if strings.Contains(line, "Date:") {
			data.Date = extractValue(line, "Date:")
		} else if strings.Contains(line, "From:") {
			data.From = extractValue(line, "From:")
		} else if strings.Contains(line, "To:") {
			data.To = extractValue(line, "To:")
		} else if strings.Contains(line, "Subject:") {
			data.Subject = extractValue(line, "Subject:")
		} else if strings.Contains(line, "Cc:") {
			data.Cc = extractValue(line, "Cc:")
		} else if strings.Contains(line, "Mime-Version:") {
			data.Mime_Version = extractValue(line, "Mime-Version:")
		} else if strings.Contains(line, "Content-Type:") {
			data.Content_Type = extractValue(line, "Content-Type:")
		} else if strings.Contains(line, "Content-Transfer-Encoding:") {
			data.Content_Transfer_Encoding = extractValue(line, "Content-Transfer-Encoding:")
		} else if strings.Contains(line, "X-From:") {
			data.X_From = extractValue(line, "X-From:")
		} else if strings.Contains(line, "X-To:") {
			data.X_To = extractValue(line, "X-To:")
		} else if strings.Contains(line, "X-cc:") {
			data.X_cc = extractValue(line, "X-cc:")
		} else if strings.Contains(line, "X-bcc:") {
			data.X_bcc = extractValue(line, "X-bcc:")
		} else if strings.Contains(line, "X-Folder:") {
			data.X_Folder = extractValue(line, "X-Folder:")
		} else if strings.Contains(line, "X-Origin:") {
			data.X_Origin = extractValue(line, "X-Origin:")
		} else if strings.Contains(line, "X-FileName:") {
			data.X_FileName = extractValue(line, "X-FileName:")
		} else {
			data.Body += line
		}
	}

	return data
}

func extractValue(line, field string) string {
	parts := strings.SplitN(line, field, 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

