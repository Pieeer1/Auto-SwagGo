package swaggo

type SwagDoc struct {
	OpenAPIVersion string                     `json:"openapi"`
	Info           Info                       `json:"info"`
	ExternalDocs   ExternalDocs               `json:"externalDocs"`
	Servers        []Server                   `json:"servers"`
	Tags           []Tag                      `json:"tags"`
	Paths          map[string]map[string]Path `json:"paths"`
	Components     Components                 `json:"components"`
}

type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

type Contact struct {
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ExternalDocs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Server struct {
	URL string `json:"url"`
}

type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type Path struct {
	Tags        []string              `json:"tags"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	OperationID string                `json:"operationId"`
	Parameters  []Parameter           `json:"parameters"`
	RequestBody *Body                 `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Schema      Schema `json:"schema"`
}

type Response struct {
	Description string `json:"description"`
	Content     map[string]Content
	Headers     map[string]Header
}

type Header struct {
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

type Body struct {
	Description string             `json:"description"`
	Content     map[string]Content `json:"content"`
	Required    bool               `json:"required"`
}

type Content struct {
	Schema Schema `json:"schema"`
}

type Schema struct {
	Type       string              `json:"type,omitempty"`
	Items      *Items              `json:"items,omitempty"`
	Ref        string              `json:"$ref,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Type        string   `json:"type,omitempty"`
	Description string   `json:"description,omitempty"`
	Format      string   `json:"format,omitempty"`
	Example     any      `json:"example,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type Items struct {
	Ref string `json:"$ref,omitempty"`
}

type Components struct {
	Schemas         map[string]Schema         `json:"schemas"`
	RequestBodies   map[string]Body           `json:"requestBodies"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
}

type SecurityScheme struct {
	Type  string `json:"type"`
	Flows Flows  `json:"flows,omitempty"`
	Name  string `json:"name,omitempty"`
	In    string `json:"in,omitempty"`
}

type Flows struct {
	Implicit Implicit `json:"implicit,omitempty"`
}

type Implicit struct {
	AuthorizationURL string            `json:"authorizationUrl"`
	Scopes           map[string]string `json:"scopes"`
}
