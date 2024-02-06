package controllers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"trucode/finalproject/models"
)

const (
	serverPort     = ":8080"
	fileServerPort = ":4080"
	batchSize      = 1000
)

var jsonFinal []string

func JSONFinal(data []string) {
	file, err := os.Create("jsonFinal.json")
	if err != nil {
		log.Fatal("failed creating json file", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(data); err != nil {
		log.Fatal("failed encoding json data: ", err)
	}
}

// It indexes an instance of IndexedData.
func IndexData(data models.Email) {
	user := os.Getenv("ZINC_FIRST_ADMIN_USER")
	password := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")

	if user == "" || password == "" {
		panic("ZINC_FIRST_ADMIN_USER and ZINC_FIRST_ADMIN_PASSWORD must be set")
	}

	auth := user + ":" + password
	bas64encoded_creds := base64.StdEncoding.EncodeToString([]byte(auth))

	zincURL := fmt.Sprintf(os.Getenv("ZINC_SEARCH_URL"), os.Getenv("INDEX_NAME"))
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshaling JSON data:", err)
	}

	jsonFinal = append(jsonFinal, string(jsonData))

	req, err := http.NewRequest("POST", zincURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error reading request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+bas64encoded_creds)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()
}

// It indexes a batch of instances of IndexedData.
func IndexBatch(batch []models.Email) {
	// Your existing batch indexing logic goes here
	for _, data := range batch {
		IndexData(data)
	}
}

func ProcessFiles(wg *sync.WaitGroup, ch chan models.Email) {
	defer wg.Done()

	myPATH := "../data/enron_mail_20110402/maildir/"
	count := 0
	batch := make([]models.Email, 0, batchSize)

	fmt.Println("Indexing...")

	err := filepath.Walk(myPATH, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking path:", err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		count++
		fmt.Println("Indexing:", path)

		sysFile, err := os.Open(path)
		if err != nil {
			log.Println("Error opening file:", err)
			return nil
		}
		defer sysFile.Close()

		lines := bufio.NewScanner(sysFile)
		indexedData := ParseData(lines, count)
		batch = append(batch, indexedData)

		if len(batch) >= batchSize {
			// Send de batch to de DB
			IndexBatch(batch)
			// Clean the batch
			batch = make([]models.Email, 0, batchSize)
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error walking the path:", err)
	}

	// Emails that not get the max len
	if len(batch) > 0 {
		IndexBatch(batch)
	}

	// JSONFinal(jsonFinal)
	close(ch)
	fmt.Println("Indexing finished!!")

}

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	filePath := "../data/enron_mail_20110402"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		log.Println("Error creating form file:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Println("Error copying file content:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writer.WriteField("key", "value")
	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost"+fileServerPort, body)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("admin", "Complexpass#123")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading HTTP response body:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(respBody))
}

// func inserData(wg *sync.WaitGroup, ch chan models.Email) {
// 	defer wg.Done()

// 	var batch []models.Email

// 	for data := range ch {
// 		batch = append(batch, data)

// 		if len(batch) >= batchSize {
// 			IndexBatch(batch)
// 			batch = nil
// 		}
// 	}

// 	// Insertar el lote final si hay registros restantes
// 	if len(batch) > 0 {
// 		IndexBatch(batch)
// 	}

// 	close(ch)

// }

// func LoadData(filePath string) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Fatalf("Error opening file: %v", err)
// 	}
// 	defer file.Close()

// 	// Create a channel to receive Email data
// 	ch := make(chan models.Email)

// 	// Use a WaitGroup to wait for all goroutines to finish
// 	var wg sync.WaitGroup

// 	// Start a goroutine to process the file
// 	wg.Add(1)
// 	go ProcessFiles(&wg, ch)

// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 	}()

// 	// Start a goroutine to insert data from channel
// 	wg.Add(1)
// 	go inserData(&wg, ch)

// 	// Wait for all goroutines to finish
// 	wg.Wait()
// }
