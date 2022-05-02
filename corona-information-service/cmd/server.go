package main

import (
	"corona-information-service/internal/handler"
	"corona-information-service/internal/handler/cases"
	"corona-information-service/internal/handler/policy"
	"corona-information-service/internal/model"
	"corona-information-service/pkg/customhttp"
	"corona-information-service/pkg/db"
	"corona-information-service/pkg/webhook"
	"corona-information-service/tools/hash"
	"log"
	"net/http"
	"os"
	"strings"
)

//Init hash secret
func init() {
	hash.Secret = []byte{0, 4, 0, 2, 2, 0, 0, 0}
}

func main() {
	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	//Initializes firestore client, if it fails then the API should still work, but without notifications/webhooks
	err := db.InitializeFirestore("./service-account.json")
	if err == nil {
		webhook.InitializeWebhooks()
		http.HandleFunc(model.NOTIFICATION_PATH, handler.NotificationHandler())
		http.HandleFunc(strings.TrimSuffix(model.NOTIFICATION_PATH, "/"), handler.NotificationHandler()) //Will be forgiving since some forget "/" at the end
	}

	// Set up handler endpoints
	http.HandleFunc(model.DEFAULT_PATH, handler.DefaultHandler())
	http.HandleFunc(model.CASE_PATH, cases.CaseHandler(customhttp.Client))
	http.HandleFunc(model.POLICY_PATH, policy.PolicyHandler(customhttp.Client))
	http.HandleFunc(model.STATUS_PATH, handler.StatusHandler(customhttp.Client))

	// Start server
	log.Println("Starting server on port " + port + " ...")

	//Will allow firestore to close the connection before terminating the server
	err = http.ListenAndServe(":"+port, nil)
	db.CloseFirestore()
	log.Fatal(err)
}
