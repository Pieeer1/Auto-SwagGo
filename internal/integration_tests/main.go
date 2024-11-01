package main

import (
	"auto-swaggo/swaggo"
	"fmt"
	"log"
	"net/http"
)

type ExampleChildrenModels struct {
	ExampleInts               []int                       `json:"example_ints" required:"true" description:"Example ints"`
	ExampleChildrenModel      ExampleChildrenModel        `json:"example_children_model" required:"true" description:"Example children model"`
	ExampleChildrenArrayModel []ExampleChildrenArrayModel `json:"example_children_array_model" required:"true" description:"Example children array model"`
}

type ExampleChildrenModel struct {
	ExampleChildrenField string `json:"example_children_field" required:"true" description:"Example children field"`
}

type ExampleChildrenArrayModel struct {
	ExampleChildrenInt int `json:"example_children_int" required:"true" description:"Example children int"`
}

type ExampleQueryStruct struct {
	ExampleQueryField    string `json:"example_query_field" name:"some custom name" required:"true" description:"Example query field"`
	ExampleIntQueryField int    `json:"example_int_query_field" required:"false" description:"Example query field"`
}

type ExampleBodyStruct struct {
	ExampleField    string `json:"example_field" required:"true" description:"Example field"`
	ExampleIntField int    `json:"example_int_field"`
}

type ExampleResponse struct {
	ExampleResponseField    string `json:"example_response_field" required:"true" description:"Example response field"`
	ExampleIntResponseField int    `json:"example_int_response_field"`
}

type ExamplePathStruct struct {
	ExamplePathField string `json:"example_path_field" name:"param" required:"true" description:"Example path field"`
}

func main() {

	authConfiguration := &swaggo.AuthenticationConfiguration{
		BasicAuth:  &swaggo.BasicAuth{},
		BearerAuth: &swaggo.BearerAuth{},
		ApiKeyAuth: &swaggo.ApiKeyAuth{
			In: "header",
		},
		OpenIdAuth: &swaggo.OpenIdAuth{
			OpenIdConnectUrl: "http://example.com",
		},
		Oauth2Auth: &swaggo.Oauth2Auth{

			Flows: swaggo.Oauth2Flows{
				Implicit: &swaggo.Oauth2Flow{
					AuthorizationUrl: "http://example.com",
					Scopes: map[string]string{
						"read": "Read access",
					},
				},
			},
		},
	}

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
		Responses: []swaggo.ResponseData{
			{
				Code: 200,
				Data: ExampleChildrenModels{
					ExampleInts: []int{1, 2, 34},
					ExampleChildrenModel: ExampleChildrenModel{
						ExampleChildrenField: "example",
					},
					ExampleChildrenArrayModel: []ExampleChildrenArrayModel{
						{
							ExampleChildrenInt: 1,
						},
					},
				},
			},
		},
		AuthenticationConfiguration: authConfiguration,
		OauthScopes:                 []string{"read"},
	})

	mux.HandleFunc("/testRouteParam/{param}", health, "v1", swaggo.RequestDetails{
		Method: "GET",
		Requests: []swaggo.RequestData{
			{
				Type: swaggo.PathSource,
				Data: ExamplePathStruct{
					ExamplePathField: "example",
				},
			},
			{
				Type: swaggo.HeaderSource,
				Data: ExamplePathStruct{
					ExamplePathField: "example",
				},
			},
		},
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
		Responses: []swaggo.ResponseData{
			{
				Code: 200,
				Data: ExampleResponse{
					ExampleResponseField:    "example",
					ExampleIntResponseField: 1,
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
