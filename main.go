package main

import (
	"auto-swaggo/swaggo"
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := swaggo.NewSwaggoMux(&swaggo.SwaggerInfo{
		Title:                   "Test API",
		Description:             "This is a test API",
		TermsOfServiceURL:       "http://example.com/terms/",
		ContactEmail:            "",
		LicenseName:             "MIT",
		LicenseURL:              "http://mit.com",
		Version:                 "1.0.0",
		ExternalDocsDescription: "Find more info here",
		ExternalDocsURL:         "http://example.com",
		Servers:                 []string{"http://localhost:8080"},
	}, "http://localhost:8080", "/api", []string{"v1", "v5"})

	mux.HandleFunc("/health", health, []string{"GET"}, "v1")
	mux.HandleFunc("/test/testing/testers", health, []string{"POST"}, "v5")
	mux.HandleFunc("/some-endpoint", health, []string{"GET", "POST"}, "")

	mux.OpenBrowser()

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", "8080"), mux)

	if err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
