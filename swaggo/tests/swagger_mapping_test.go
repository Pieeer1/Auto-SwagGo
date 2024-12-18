package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/Pieeer1/Auto-SwagGo/swaggo"
)

// good test because it covers a large amount of the types possible
type TestChildrenModels struct {
	ExampleString             string
	ExampleInts               []int
	ExampleChildrenModel      TestChildrenModel
	ExampleChildrenArrayModel []TestChildrenArrayModel
}

type TestChildrenModel struct {
	ExampleChildrenField string
}

type TestChildrenArrayModel struct {
	ExampleChildrenInt int
}

type TimeTestChildrenModel struct {
	ExampleTime    time.Time
	ExampleTimePtr *time.Time
}

func TestSwaggerMapWithAuth(t *testing.T) {
	swaggoMux := swaggo.NewSwaggoMux(&swaggo.SwaggerInfo{
		Title: "Test",
	}, "http://test:8080", "/api", []string{"v1"})

	swaggoMux.HandleFunc("/test", nil, "v1", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
		AuthenticationConfiguration: &swaggo.AuthenticationConfiguration{
			Oauth2Auth: &swaggo.Oauth2Auth{
				Flows: swaggo.Oauth2Flows{
					Implicit: &swaggo.Oauth2Flow{
						AuthorizationUrl: "http://example.com",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				},
			},
		},
		OauthScopes: []string{"read"},
	})

	doc, err := swaggoMux.MapDoc("")

	if err != nil {
		t.Fatal(err)
	}

	if doc.Paths["/api/v1/test"]["get"].Security[0]["oauth2"] == nil {
		t.Errorf("Expected oauth2, got %s", doc.Paths["/api/v1/test"]["get"].Security[0]["oauth2"])
	}

	if doc.Components.SecuritySchemes["oauth2"].Flows.Implicit == nil {
		t.Errorf("Expected implicit oauth2 to exist")
	}

	if doc.Components.SecuritySchemes["oauth2"].Flows.Implicit.AuthorizationURL != "http://example.com" {
		t.Errorf("Expected http://example.com, got %s", doc.Components.SecuritySchemes["oauth2"].Flows.Implicit.AuthorizationURL)
	}

	if doc.Components.SecuritySchemes["oauth2"].Flows.Implicit.Scopes["read"] != "Read access" {
		t.Errorf("Expected Read access, got %s", doc.Components.SecuritySchemes["oauth2"].Flows.Implicit.Scopes["read"])
	}

	if doc.Components.SecuritySchemes["oauth2"].Flows.Implicit.Scopes["write"] != "Write access" {
		t.Errorf("Expected Write access, got %s", doc.Components.SecuritySchemes["oauth2"].Flows.Implicit.Scopes["write"])
	}

	if doc.Components.SecuritySchemes["oauth2"].Flows.Password != nil {
		t.Errorf("Expected nil, got %s", doc.Components.SecuritySchemes["oauth2"].Flows.Password)
	}

}

func TestBaseSwaggerMap(t *testing.T) {

	swaggoMux := swaggo.NewSwaggoMux(&swaggo.SwaggerInfo{
		Title: "Test",
	}, "http://test:8080", "/api", []string{"v1"})

	swaggoMux.HandleFunc("/test", nil, "v1", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
		Requests: []swaggo.RequestData{
			{
				Type:        swaggo.QuerySource,
				Description: "Test",
				Required:    true,
				ContentType: []string{"application/json"},
			},
		},
		Responses: []swaggo.ResponseData{
			{
				Code: 200,
				Data: TestChildrenModels{
					ExampleString: "example",
					ExampleInts:   []int{1, 2, 34},
					ExampleChildrenModel: TestChildrenModel{
						ExampleChildrenField: "second example",
					},
					ExampleChildrenArrayModel: []TestChildrenArrayModel{
						{
							ExampleChildrenInt: 1,
						},
					},
				},
			},
		},
	})

	doc, err := swaggoMux.MapDoc("")

	if err != nil {
		t.Fatal(err)
	}

	if doc.Info.Title != "Test" {
		t.Errorf("Expected Test , got %s", doc.Info.Title)
	}

	if doc.Paths["/api/v1/test"]["get"].Summary != "Test" {
		t.Errorf("Expected Test, got %s", doc.Paths["/api/v1/test"]["get"].Summary)
	}

	if doc.Paths["/api/v1/test"]["get"].Description != "Test" {
		t.Errorf("Expected Test, got %s", doc.Paths["/api/v1/test"]["get"].Description)
	}

	if doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Ref != "#/components/schemas/TestChildrenModels" {
		t.Errorf("Expected object, got %s", doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Type != "string" {
		t.Errorf("Expected string, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Example != "example" {
		t.Errorf("Expected string, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Type != "array" {
		t.Errorf("Expected array, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Items.Type != "integer" {
		t.Errorf("Expected integer, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Items.Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleChildrenModel"].Type != "object" {
		t.Errorf("Expected object, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleChildrenModel"].Properties["ExampleChildrenField"].Type != "string" {
		t.Errorf("Expected string, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Type)
	}

}

func TestSwaggerMappingArray(t *testing.T) {

	swaggoMux := swaggo.NewSwaggoMux(&swaggo.SwaggerInfo{
		Title: "Test",
	}, "http://test:8080", "/api", []string{"v1"})

	swaggoMux.HandleFunc("/test", nil, "v1", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
		Requests: []swaggo.RequestData{
			{
				Type: swaggo.BodySource,
				Data: []TestChildrenModels{
					{
						ExampleString: "example",
						ExampleInts:   []int{1, 2, 34},
						ExampleChildrenModel: TestChildrenModel{
							ExampleChildrenField: "second example",
						},
						ExampleChildrenArrayModel: []TestChildrenArrayModel{
							{
								ExampleChildrenInt: 1,
							},
						},
					},
				},
				Description: "Test",
				Required:    true,
				ContentType: []string{"application/json"},
			},
		},
		Responses: []swaggo.ResponseData{
			{
				Code: 200,
				Data: []TestChildrenModels{ // TODO - FIX BUG WHERE EMPTY ARRAY CHILDREN WILL PANIC
					{
						ExampleString: "example",
						ExampleInts:   []int{1, 2, 34},
						ExampleChildrenModel: TestChildrenModel{
							ExampleChildrenField: "second example",
						},
						ExampleChildrenArrayModel: []TestChildrenArrayModel{
							{
								ExampleChildrenInt: 1,
							},
						},
					},
				},
			},
		},
	})

	doc, err := swaggoMux.MapDoc("")

	if err != nil {
		t.Fatal(err)
	}

	test := doc.Paths["/api/v1/test"]["get"].Responses["200"]

	fmt.Printf("%+v\n", test)

	if doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Type != "array" {
		t.Errorf("Expected array, got %s", doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Type)
	}

	if doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Items.Ref != "#/components/schemas/TestChildrenModels" {
		t.Errorf("Expected valdi schema item ref, got %s", doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Items.Ref)
	}

	if doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Type != "array" {
		t.Errorf("Expected array, got %s", doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Type)
	}

	if doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Items.Ref != "#/components/schemas/TestChildrenModels" {
		t.Errorf("Expected valdi schema item ref, got %s", doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Items.Ref)
	}

}

func TestSwaggerMappingEmptyChildArray(t *testing.T) {
	swaggoMux := swaggo.NewSwaggoMux(&swaggo.SwaggerInfo{
		Title: "Test",
	}, "http://test:8080", "/api", []string{"v1"})

	swaggoMux.HandleFunc("/test", nil, "v1", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
		Requests: []swaggo.RequestData{
			{
				Type:        swaggo.BodySource,
				Data:        []TestChildrenModels{}, // point of this test. Need to validate empty array's can be passed in request body
				Description: "Test",
				Required:    true,
				ContentType: []string{"application/json"},
			},
		},
		Responses: []swaggo.ResponseData{
			{
				Code: 200,
				Data: []TestChildrenModels{}, // as well as empty array's in response
			},
		},
	})

	doc, err := swaggoMux.MapDoc("")

	if err != nil {
		t.Fatal(err)
	}

	test := doc.Paths["/api/v1/test"]["get"].Responses["200"]

	fmt.Printf("%+v\n", test)

	if doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Type != "array" {
		t.Errorf("Expected array, got %s", doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Type)
	}

	if doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Items.Ref != "#/components/schemas/TestChildrenModels" {
		t.Errorf("Expected valdi schema item ref, got %s", doc.Paths["/api/v1/test"]["get"].Responses["200"].Content["application/json"].Schema.Items.Ref)
	}

	if doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Type != "array" {
		t.Errorf("Expected array, got %s", doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Type)
	}

	if doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Items.Ref != "#/components/schemas/TestChildrenModels" {
		t.Errorf("Expected valdi schema item ref, got %s", doc.Paths["/api/v1/test"]["get"].RequestBody.Content["application/json"].Schema.Items.Ref)
	}
}

func TestSwaggerMappingVersioning(t *testing.T) {
	swaggoMux := swaggo.NewSwaggoMux(&swaggo.SwaggerInfo{
		Title: "Test",
	}, "http://test:8080", "/api", []string{"v1", "v2"})

	swaggoMux.HandleFunc("/test", nil, "", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
	})

	swaggoMux.HandleFunc("/test", nil, "v1", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
		Requests: []swaggo.RequestData{
			{
				Type: swaggo.BodySource,
				Data: TimeTestChildrenModel{
					ExampleTime:    time.Now(),
					ExampleTimePtr: nil,
				},
			},
		},
		Responses: []swaggo.ResponseData{
			{
				Code: 200,
				Data: TestChildrenModels{
					ExampleString: "example",
				},
			},
		},
	})

	swaggoMux.HandleFunc("/test", nil, "v2", swaggo.RequestDetails{
		Method:      "GET",
		Summary:     "Test",
		Description: "Test",
	})

	doc, err := swaggoMux.MapDoc("")

	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Paths) != 3 {
		t.Errorf("Expected 3 paths, got %d", len(doc.Paths))
	}

	if doc.Paths["/api/test"]["get"].Summary != "Test" {
		t.Errorf("Expected Test, got %s", doc.Paths["/api/test"]["get"].Summary)
	}

	if doc.Paths["/api/v1/test"]["get"].Summary != "Test" {
		t.Errorf("Expected Test, got %s", doc.Paths["/api/v1/test"]["get"].Summary)
	}

	if doc.Paths["/api/v2/test"]["get"].Summary != "Test" {
		t.Errorf("Expected Test, got %s", doc.Paths["/api/v2/test"]["get"].Summary)
	}

	if len(doc.Components.Schemas) != 2 {
		t.Errorf("Expected 2 schemas, got %d", len(doc.Components.Schemas))
	}

	docv1, err := swaggoMux.MapDoc("v1")

	if err != nil {
		t.Fatal(err)
	}

	if len(docv1.Paths) != 1 {
		t.Errorf("Expected 1 paths, got %d", len(docv1.Paths))
	}

	if docv1.Paths["/api/v1/test"]["get"].Summary != "Test" {
		t.Errorf("Expected Test, got %s", docv1.Paths["/api/v1/test"]["get"].Summary)
	}

	if len(docv1.Components.Schemas) != 2 {
		t.Errorf("Expected 2 schemas, got %d", len(doc.Components.Schemas))
	}

	docv2, err := swaggoMux.MapDoc("v2")

	if err != nil {
		t.Fatal(err)
	}

	if len(docv2.Paths) != 1 {
		t.Errorf("Expected 1 paths, got %d", len(docv2.Paths))
	}

	if docv2.Paths["/api/v2/test"]["get"].Summary != "Test" {
		t.Errorf("Expected Test, got %s", docv2.Paths["/api/v2/test"]["get"].Summary)
	}

	if len(docv2.Components.Schemas) != 0 {
		t.Errorf("Expected 0 schemas, got %d", len(docv2.Components.Schemas))
	}

}
