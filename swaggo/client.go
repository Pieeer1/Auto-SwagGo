package swaggo

import (
	"auto-swaggo/internal/ext"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

var IGNORED_TAGS = []string{"swagger", "openapi.json"}

type SwaggoMux struct {
	mux         *http.ServeMux
	swaggerInfo *SwaggerInfo
	baseUri     string
	prefix      string
	versions    []string
	routes      []Route
	mu          sync.RWMutex
}

func NewSwaggoMux(swaggerInfo *SwaggerInfo, baseUri, prefix string, versions []string) *SwaggoMux {
	client := &SwaggoMux{
		routes:      make([]Route, 0),
		swaggerInfo: swaggerInfo,
		baseUri:     baseUri,
		versions:    versions,
		prefix:      prefix,
		mux:         http.NewServeMux(),
		mu:          sync.RWMutex{},
	}

	if len(versions) == 0 {
		client.HandleFunc("/swagger/index.html", client.swagger, "", RequestDetails{Method: "GET"})
		client.HandleFunc("/openapi.json", client.swaggerJson, "", RequestDetails{Method: "GET"})
	}

	for _, version := range versions {
		client.HandleFunc("/swagger/index.html", client.swagger, version, RequestDetails{Method: "GET"})
		client.HandleFunc("/openapi.json", client.swaggerJson, version, RequestDetails{Method: "GET"})
	}

	return client
}

func (m *SwaggoMux) Handle(path string, handler http.Handler, version string, requestDetails ...RequestDetails) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var fullPath string

	if version == "" {
		fullPath = fmt.Sprintf("%s%s", m.prefix, path)
	} else {
		fullPath = fmt.Sprintf("%s/%s%s", m.prefix, version, path)
	}

	m.routes = append(m.routes, Route{Path: fullPath, Handler: handler, Prefix: m.prefix, Version: version, RequestDetails: requestDetails})
	m.mux.Handle(fullPath, m.defaultMiddleware(handler, requestDetails))

}

func (m *SwaggoMux) defaultMiddleware(handler http.Handler, requestDetails []RequestDetails) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		methods := ext.SliceMap(requestDetails, func(rd RequestDetails) string {
			return rd.Method
		})

		if !ext.Contains(methods, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		handler.ServeHTTP(w, r)
	})
}

func (m *SwaggoMux) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request), version string, requestDetails ...RequestDetails) {
	m.Handle(path, http.HandlerFunc(handler), version, requestDetails...)
}

func (m *SwaggoMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (c *SwaggoMux) OpenBrowser() {
	var endpoint string

	if len(c.versions) == 0 {
		endpoint = fmt.Sprintf("%s%s/swagger/index.html", c.baseUri, c.prefix)
	} else {
		endpoint = fmt.Sprintf("%s%s/%s/swagger/index.html", c.baseUri, c.prefix, c.versions[0])
	}

	openBrowser(endpoint)
}

func (c *SwaggoMux) swaggerJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mappedDoc, err := c.mapDoc()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	docRb, err := json.Marshal(mappedDoc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(docRb)
}

func (c *SwaggoMux) swagger(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var endpoint string

	if len(c.versions) == 0 {
		endpoint = fmt.Sprintf("%s%s/openapi.json", c.baseUri, c.prefix)
	} else {
		endpoint = fmt.Sprintf("%s%s/%s/openapi.json", c.baseUri, c.prefix, c.versions[0])
	}

	w.Write([]byte(fmt.Sprintf(
		`
		<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<meta name="description" content="SwaggerUI" />
				<title>SwaggerUI</title>
				<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
			</head>
			<body>
			<div id="swagger-ui"></div>
			<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
				<script>
					window.onload = () => {
						window.ui = SwaggerUIBundle({
						url: '%s', 
						dom_id: '#swagger-ui',
						});
					};
				</script>
			</body>
		</html>
		`, endpoint)))
}

func (c *SwaggoMux) mapDoc() (*SwagDoc, error) {

	tagNames := ext.SliceMap(c.routes, func(route Route) string {
		return route.GetPathWithoutPrefixAndVersion()
	})

	tags := ext.SliceMap(ext.Where(ext.Distinct(tagNames), func(tagName string) bool {
		return !ext.Contains(IGNORED_TAGS, tagName)
	}), func(tagName string) Tag {
		return Tag{Name: tagName, Description: fmt.Sprintf("Operations for %s", tagName)}
	})

	paths := make(map[string]map[string]Path)

	for _, route := range ext.Where(c.routes, func(route Route) bool {
		return !ext.Contains(IGNORED_TAGS, route.GetPathWithoutPrefixAndVersion())
	}) {

		paths[route.Path] = make(map[string]Path)

		tagName := route.GetPathWithoutPrefixAndVersion()

		for _, rd := range route.RequestDetails {

			queryRequests := ext.Where(rd.Requests, func(rd RequestData) bool {
				return rd.Type == QuerySource
			})

			parameters := make([]Parameter, 0)

			for _, qr := range queryRequests {
				t, v, err := rawReflect(qr.Data)

				if err != nil {
					return nil, err
				}

				for i := 0; i < t.NumField(); i++ {
					field := t.Field(i)
					value := v.Field(i)

					var fName string

					if field.Tag.Get("name") != "" {
						fName = field.Tag.Get("name")
					} else {
						fName = field.Name
					}

					parameters = append(parameters, Parameter{
						Name:        fName,
						In:          "query",
						Description: field.Tag.Get("description"),
						Required:    field.Tag.Get("required") == "true",
						Schema:      Schema{Type: parseGOTypeToSwaggerType(value.Kind())},
					})
				}
			}

			bodyRequests := ext.Where(rd.Requests, func(rd RequestData) bool {
				return rd.Type == BodySource
			})

			var body *Body

			if len(bodyRequests) > 1 {
				return nil, fmt.Errorf("only one body request is allowed")
			} else if len(bodyRequests) == 1 {
				for _, br := range bodyRequests {
					body = &Body{
						Content:     map[string]Content{},
						Description: br.Description,
						Required:    br.Required,
					}

					if len(br.ContentType) == 0 {
						br.ContentType = []string{"application/json"} // default to application/json if no type is given
					}

					splitTypeName := strings.Split(reflect.TypeOf(br.Data).String(), ".")

					for _, contentType := range br.ContentType {
						body.Content[contentType] = Content{
							Schema: Schema{
								Ref: fmt.Sprintf("#/components/schemas/%s", splitTypeName[len(splitTypeName)-1]),
							},
						}
					}

				}
			}

			responses := map[string]Response{}

			if len(rd.Responses) > 0 {
				for _, res := range rd.Responses {
					splitTypeName := strings.Split(reflect.TypeOf(res.Data).String(), ".")

					content := map[string]Content{}

					if len(res.ContentType) == 0 {
						res.ContentType = []string{"application/json"} // default to application/json if no type is given
					}

					for _, contentType := range res.ContentType {
						content[contentType] = Content{
							Schema: Schema{
								Ref: fmt.Sprintf("#/components/schemas/%s", splitTypeName[len(splitTypeName)-1]),
							},
						}
					}

					headerMap := make(map[string]Header)

					for header, value := range res.Headers {
						headerMap[header] = Header{Schema: Schema{Type: parseGOTypeToSwaggerType(reflect.TypeOf(value).Kind())}}
					}

					responses[fmt.Sprintf("%d", res.Code)] = Response{
						Headers:     headerMap,
						Description: fmt.Sprintf("%d response", res.Code),
						Content:     content,
					}
				}
			} else {
				responses[""] = Response{
					Description: "",
				}
			}

			paths[route.Path][strings.ToLower(rd.Method)] = Path{
				Tags:        []string{tagName},
				Summary:     rd.Summary,
				Description: rd.Description,
				OperationID: fmt.Sprintf("%s-%s", rd.Method, route.Path),
				Parameters:  parameters,
				RequestBody: body,
				Responses:   responses,
				Security:    []map[string][]string{},
			}
		}

	}

	schemas := make(map[string]Schema)

	distinctRequestTypes := ext.DistinctBy(ext.FlattenMap(c.routes, func(route Route) []RequestData {
		return ext.FlattenMap(route.RequestDetails, func(requestDetails RequestDetails) []RequestData {
			return requestDetails.Requests
		})
	}), func(reqBody RequestData) string {
		return reflect.TypeOf(reqBody.Data).String()
	})

	distinctResponseTypes := ext.DistinctBy(ext.FlattenMap(c.routes, func(route Route) []ResponseData {
		return ext.FlattenMap(route.RequestDetails, func(requestDetails RequestDetails) []ResponseData {
			return requestDetails.Responses
		})
	}), func(reqBody ResponseData) string {
		return reflect.TypeOf(reqBody.Data).String()
	})

	distinctTypes := append(ext.SliceMap(distinctRequestTypes, func(req RequestData) any {
		return req.Data
	}), ext.SliceMap(distinctResponseTypes, func(res ResponseData) any {
		return res.Data
	})...)

	for _, data := range distinctTypes {
		t, v, err := rawReflect(data)

		if err != nil {
			return nil, err
		}

		properties := make(map[string]Property)

		requiredProperties := make([]string, 0)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			var fName string

			if field.Tag.Get("name") != "" {
				fName = field.Tag.Get("name")
			} else {
				fName = field.Name
			}

			if field.Tag.Get("required") == "true" {
				requiredProperties = append(requiredProperties, fName)
			}

			properties[fName] = Property{
				Type:        parseGOTypeToSwaggerType(value.Kind()),
				Description: field.Tag.Get("description"),
				Example:     autoType(value.Kind(), value),
			}
		}

		var splitSchemaName = strings.Split(t.String(), ".")

		schemas[splitSchemaName[len(splitSchemaName)-1]] = Schema{
			Type:       "object",
			Properties: properties,
			Required:   requiredProperties,
		}
	}

	doc := &SwagDoc{
		OpenAPIVersion: "3.0.2",
		Info: Info{
			Title:          c.swaggerInfo.Title,
			Description:    c.swaggerInfo.Description,
			TermsOfService: c.swaggerInfo.TermsOfServiceURL,
			Version:        c.swaggerInfo.Version,
			Contact: Contact{
				Email: c.swaggerInfo.ContactEmail,
			},
			License: License{
				Name: c.swaggerInfo.LicenseName,
				URL:  c.swaggerInfo.LicenseURL,
			},
		},
		ExternalDocs: ExternalDocs{
			Description: c.swaggerInfo.ExternalDocsDescription,
			URL:         c.swaggerInfo.ExternalDocsURL,
		},
		Servers: ext.SliceMap(c.swaggerInfo.Servers, func(serverUri string) Server {
			return Server{URL: serverUri}
		}),
		Tags:  tags,
		Paths: paths,
		Components: Components{
			Schemas: schemas,
		},
	}

	return doc, nil
}

func rawReflect(data any) (reflect.Type, reflect.Value, error) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, reflect.Value{}, fmt.Errorf("data must be a struct or a pointer to a struct")
	}

	return t, v, nil
}

func parseGOTypeToSwaggerType(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Map, reflect.Pointer, reflect.Interface, reflect.Struct, reflect.UnsafePointer:
		return "object"
	default:
		return "object"
	}
}

func autoType(kind reflect.Kind, value reflect.Value) any {
	switch kind {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Bool:
		return value.Bool()
	case reflect.Array, reflect.Slice:
		return value.Slice(0, value.Len())
	case reflect.Map, reflect.Pointer, reflect.Interface, reflect.Struct, reflect.UnsafePointer:
		return value.Interface()
	default:
		return value.Interface()
	}
}
