package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"trucode/search/utils"
)

type JsonBody struct {
	SearchTerm string `json:"searchTerm"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	var body JsonBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.JsonWriter(w, http.StatusBadRequest, "invalid request body") // Status Error 400
		return
	}

	data, err := GetEmails(body.SearchTerm)
	if err != nil {
		log.Println("error", err)
		utils.JsonWriter(w, http.StatusInternalServerError, "error getting emails") // Status Error 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
