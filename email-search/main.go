package main

import (
	"log"
	"net/http"
	"os"
	"trucode/search/api"
	"trucode/search/env"
)

func main() {
	env.LoadVars()

	port := os.Getenv("SERVER_PORT")

	server := api.SetupRouter()

	log.Printf("Server running on port %s", port)

	log.Fatal(http.ListenAndServe(port, server))

}
