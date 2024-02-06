package main

import (
	"net/http"
	"sync"
	"trucode/finalproject/api"
	"trucode/finalproject/controllers"
	"trucode/finalproject/env"
	"trucode/finalproject/models"
)

const (
	serverPort = ":8080"
)

func main() {
	env.LoadVars()
	router := api.SetupRouter()
	router.Post("/enron", controllers.HandleFileUpload)

	var wg sync.WaitGroup
	ch := make(chan models.Email)
	wg.Add(1)
	go controllers.ProcessFiles(&wg, ch)

	http.ListenAndServe(serverPort, router)
	wg.Wait()
}
