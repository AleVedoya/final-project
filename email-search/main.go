package main

import (
	"encoding/json"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		message := map[string]string{"message": "El puerto está funcionando"}
		json.NewEncoder(w).Encode(message)
	})
	log.Fatal(http.ListenAndServe(port, server))

}
