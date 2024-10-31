package main

import (
	"auto-swaggo/swaggo"
	"fmt"
	"log"
	"net/http"
)

type ExampleQueryStruct struct {
	ExampleQueryField    string `json:"example_query_field" name:"some custom name" required:"true" description:"Example query field"`
	ExampleIntQueryField int    `json:"example_int_query_field" required:"false" description:"Example query field"`
}

type ExampleBodyStruct struct {
	ExampleField    string `json:"example_field" required:"true" description:"Example field"`
	ExampleIntField int    `json:"example_int_field"`
}

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

	mux.HandleFunc("/health", health, "v1", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Health Check",
		Description: "Check the health of the API",
	})

	mux.HandleFunc("/test/testing/testers", health, "v5", swaggo.RequestDetails{
		Method: "POST",
	})

	mux.HandleFunc("/some-endpoint", health, "", swaggo.RequestDetails{
		Method: "GET",
		Requests: []swaggo.RequestData{
			{
				Type: swaggo.QuerySource,
				Data: ExampleQueryStruct{
					ExampleQueryField:    "example",
					ExampleIntQueryField: 1,
				},
			},
		},
	}, swaggo.RequestDetails{
		Method: "POST",
		Requests: []swaggo.RequestData{
			{
				Type: swaggo.BodySource,
				Data: ExampleBodyStruct{
					ExampleField:    "example",
					ExampleIntField: 1,
				},
			},
		},
	})

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
