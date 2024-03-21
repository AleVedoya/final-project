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

// "../data/enron_mail_20110402/maildir/bailey-s"
// debra.bailey@enron.com (SI)
// Credit/Legal (SI)
// <23746398.1075841959452.JavaMail.evans@thyme> (NO)
// Gentlemen (SI)
// cerrar (NO)