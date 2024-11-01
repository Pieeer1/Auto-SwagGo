package tests

import (
	"auto-swaggo/internal/ext"
	"auto-swaggo/swaggo"
	"testing"
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

	doc, err := swaggoMux.MapDoc()

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

	if ext.SequenceEqual(doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Example.([]interface{}), []any{"integer"}) {
		t.Errorf("Expected string, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleChildrenModel"].Type != "object" {
		t.Errorf("Expected object, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleInts"].Type)
	}

	if doc.Components.Schemas["TestChildrenModels"].Properties["ExampleChildrenModel"].Example.(TestChildrenModel).ExampleChildrenField != "second example" {
		t.Errorf("Expected string, got %s", doc.Components.Schemas["TestChildrenModels"].Properties["ExampleString"].Type)
	}

}
