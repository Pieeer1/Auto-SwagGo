package swaggo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Pieeer1/Auto-SwagGo/internal/ext"
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

	client.HandleFunc("/swagger/index.html", client.swagger, "", RequestDetails{Method: "GET"})
	client.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		client.swaggerJson(w, "")
	}, "", RequestDetails{Method: "GET"})

	for _, version := range versions {
		client.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
			client.swaggerJson(w, version)
		}, version, RequestDetails{Method: "GET"})
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
	openBrowser(fmt.Sprintf("%s%s/swagger/index.html", c.baseUri, c.prefix))
}

func (c *SwaggoMux) swaggerJson(w http.ResponseWriter, version string) {
	w.Header().Set("Content-Type", "application/json")

	mappedDoc, err := c.MapDoc(version)

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

type VersionedUrlSwagger struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func (c *SwaggoMux) swagger(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var swaggerHtml string

	_, fileName, _, _ := runtime.Caller(0)

	filePath := filepath.Dir(fileName)

	baseHtml, err := os.ReadFile(fmt.Sprintf("%s/html/index.html", filePath))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error reading swagger html file: %s", err.Error())))
		return
	}

	if len(c.versions) == 0 {
		endpoint := fmt.Sprintf("%s%s/openapi.json", c.baseUri, c.prefix)
		swaggerHtml = strings.Replace(string(baseHtml), "{{page_url}}", endpoint, 1)
	} else {

		versionedUrls := ext.SliceMap(c.versions, func(version string) VersionedUrlSwagger {
			return VersionedUrlSwagger{
				Url:  fmt.Sprintf("%s%s/%s/openapi.json", c.baseUri, c.prefix, version),
				Name: version,
			}
		})

		versionedUrls = ext.Push(versionedUrls, VersionedUrlSwagger{
			Url:  fmt.Sprintf("%s%s/openapi.json", c.baseUri, c.prefix),
			Name: "All",
		})

		versionedUrlJsonString, err := json.Marshal(versionedUrls)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error marshalling versioned urls: %s", err.Error())))
		}

		swaggerHtml = strings.Replace(string(baseHtml), "{{is_versioned}}", "true", 1)
		swaggerHtml = strings.Replace(swaggerHtml, "{{versioned_urls}}", string(versionedUrlJsonString), 1)
		swaggerHtml = strings.Replace(swaggerHtml, "{{default_name}}", c.versions[0], 1)
	}

	w.Write([]byte(swaggerHtml))
}

func (c *SwaggoMux) MapDoc(version string) (*SwagDoc, error) {

	tagNames := ext.SliceMap(ext.Where(c.routes, func(route Route) bool {
		return (version == "" || route.Version == version)
	}), func(route Route) string {
		return route.GetPathWithoutPrefixAndVersion()
	})

	tags := ext.SliceMap(ext.Where(ext.Distinct(tagNames), func(tagName string) bool {
		return !ext.Contains(IGNORED_TAGS, tagName)
	}), func(tagName string) Tag {
		return Tag{Name: tagName, Description: fmt.Sprintf("Operations for %s", tagName)}
	})

	paths, err := c.getPaths(version)

	if err != nil {
		return nil, err
	}

	schemas, err := c.getSchemas(version)

	if err != nil {
		return nil, err
	}

	requestBodies, err := c.getRequestBodies(version)

	if err != nil {
		return nil, err
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
			Schemas:         schemas,
			RequestBodies:   requestBodies,
			SecuritySchemes: c.getSecuritySchemas(),
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

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()

		if v.Len() == 0 {
			v = reflect.MakeSlice(reflect.SliceOf(t), 1, 1)
		}

		v = v.Index(0)
	}

	kind := t.Kind()

	if !(kind == reflect.Struct || kind == reflect.Array || kind == reflect.Slice) {
		return nil, reflect.Value{}, fmt.Errorf("data must be a struct or a pointer to a struct, or an array. Got %s", kind.String())
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

func isArray(data any) bool {
	t := reflect.TypeOf(data)
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

func isByteArray(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Uint8 // uint8 is byte in reflect package
}

func isTime(field reflect.StructField) bool {
	return field.Type == reflect.TypeOf(time.Time{}) || field.Type == reflect.TypeOf(&time.Time{})
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
		arr := value.Slice(0, value.Len())
		res := make([]any, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			res[i] = autoType(arr.Index(i).Kind(), arr.Index(i))
		}
		return res
	case reflect.Map, reflect.Pointer, reflect.Interface, reflect.Struct, reflect.UnsafePointer:
		return value.Interface()
	default:
		return value.Interface()
	}
}

func (c *SwaggoMux) getRequestBodies(version string) (map[string]Body, error) {

	requestBodies := make(map[string]Body)

	distinctRequestBodies := ext.DistinctBy(ext.Where(ext.FlattenMap(ext.FlattenMap(ext.Where(c.routes, func(route Route) bool {
		return (version == "" || route.Version == version)
	}), func(route Route) []RequestDetails {
		return route.RequestDetails
	}), func(requestDetails RequestDetails) []RequestData {
		return requestDetails.Requests
	}), func(requestData RequestData) bool {
		return requestData.Type == BodySource && requestData.Data != nil
	}), func(requestData RequestData) string {
		return reflect.TypeOf(requestData.Data).String()
	})

	for _, reqBody := range distinctRequestBodies {

		content := map[string]Content{}

		splitSchemaName := strings.Split(reflect.TypeOf(reqBody.Data).String(), ".")
		friendlyName := splitSchemaName[len(splitSchemaName)-1]

		if len(reqBody.ContentType) == 0 {
			reqBody.ContentType = []string{"application/json"} // default to application/json if no type is given
		}

		for _, contentType := range reqBody.ContentType {
			content[contentType] = Content{
				Schema: Schema{
					Ref: fmt.Sprintf("#/components/schemas/%s", friendlyName),
				},
			}
		}

		requestBodies[friendlyName] = Body{
			Content:     content,
			Description: reqBody.Description,
			Required:    reqBody.Required,
		}

	}

	return requestBodies, nil
}

func (c *SwaggoMux) getSchemas(version string) (map[string]Schema, error) {
	schemas := make(map[string]Schema)

	distinctRequestTypes := ext.DistinctBy(ext.FlattenMap(ext.Where(c.routes, func(route Route) bool {
		return (version == "" || route.Version == version)
	}), func(route Route) []RequestData {
		return ext.FlattenMap(route.RequestDetails, func(requestDetails RequestDetails) []RequestData {
			return requestDetails.Requests
		})
	}), func(reqBody RequestData) string {
		if reqBody.Data == nil {
			return ""
		}
		return reflect.TypeOf(reqBody.Data).String()
	})

	distinctResponseTypes := ext.DistinctBy(ext.FlattenMap(ext.Where(c.routes, func(route Route) bool {
		return (version == "" || route.Version == version)
	}), func(route Route) []ResponseData {
		return ext.FlattenMap(route.RequestDetails, func(requestDetails RequestDetails) []ResponseData {
			return requestDetails.Responses
		})
	}), func(reqBody ResponseData) string {
		if reqBody.Data == nil {
			return ""
		}
		return reflect.TypeOf(reqBody.Data).String()
	})

	distinctTypes := append(ext.SliceMap(distinctRequestTypes, func(req RequestData) any {
		return req.Data
	}), ext.SliceMap(distinctResponseTypes, func(res ResponseData) any {
		return res.Data
	})...)

	for _, data := range distinctTypes {

		if data == nil {
			continue
		}

		t, v, err := rawReflect(data)

		if err != nil {
			return nil, err
		}

		schema, err := mapChildPropertiesToSchema(t, v)

		if err != nil {
			return nil, err
		}

		var splitSchemaName = strings.Split(t.String(), ".")

		schemas[splitSchemaName[len(splitSchemaName)-1]] = schema
	}
	return schemas, nil
}

func mapChildPropertiesToSchema(t reflect.Type, v reflect.Value) (Schema, error) {

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

		swagType := parseGOTypeToSwaggerType(value.Kind())

		if swagType == "array" && isByteArray(field) {
			properties[fName] = Property{
				Type:        "string",
				Format:      "binary",
				Description: field.Tag.Get("description"),
			}
			continue
		}

		if swagType == "array" {

			arrayType := value.Type().Elem()

			var childItemValue Schema

			parsedType := parseGOTypeToSwaggerType(arrayType.Kind())

			if parsedType == "array" {

				if value.Len() == 0 {
					value = reflect.MakeSlice(reflect.SliceOf(arrayType), 1, 1)
				}

				arrayValue := value.Index(0)

				newT, newV, err := rawReflect(autoType(arrayValue.Kind(), arrayValue))

				if err != nil {
					return Schema{}, err
				}

				childItemValue, err = mapChildPropertiesToSchema(newT, newV)

				if err != nil {
					return Schema{}, err
				}
			} else if parsedType == "object" {
				newT, newV, err := rawReflect(autoType(arrayType.Kind(), value))

				if err != nil {
					return Schema{}, err
				}

				childItemValue, err = mapChildPropertiesToSchema(newT, newV)

				if err != nil {
					return Schema{}, err
				}
			} else {
				childItemValue = Schema{
					Type: parsedType,
				}
			}

			properties[fName] = Property{
				Type:        "array",
				Items:       &childItemValue,
				Description: field.Tag.Get("description"),
			}
			continue
		} else if swagType == "object" && isTime(field) {
			properties[fName] = Property{
				Type:        "string",
				Format:      "date-time",
				Description: field.Tag.Get("description"),
			}
			continue
		} else if swagType == "object" {

			newT, newV, err := rawReflect(autoType(value.Kind(), value))

			if err != nil {
				return Schema{}, err
			}

			childItemValue, err := mapChildPropertiesToSchema(newT, newV)

			if err != nil {
				return Schema{}, err
			}

			properties[fName] = Property{
				Type:        "object",
				Properties:  childItemValue.Properties,
				Description: field.Tag.Get("description"),
			}
			continue
		}

		properties[fName] = Property{
			Type:        swagType,
			Description: field.Tag.Get("description"),
			Example:     autoType(value.Kind(), value),
		}
	}

	return Schema{
		Type:       "object",
		Properties: properties,
		Required:   requiredProperties,
	}, nil
}

func (c *SwaggoMux) getPaths(version string) (map[string]map[string]Path, error) {
	paths := make(map[string]map[string]Path)

	for _, route := range ext.Where(c.routes, func(route Route) bool {
		return !ext.Contains(IGNORED_TAGS, route.GetPathWithoutPrefixAndVersion()) && (version == "" || route.Version == version)
	}) {

		paths[route.Path] = make(map[string]Path)

		tagName := route.GetPathWithoutPrefixAndVersion()

		for _, rd := range route.RequestDetails {

			parameterRequests := ext.Where(rd.Requests, func(rd RequestData) bool {
				return ext.Contains([]RequestDataSource{QuerySource, PathSource, HeaderSource}, rd.Type)
			})

			parameters := make([]Parameter, 0)

			for _, qr := range parameterRequests {

				if qr.Data == nil {
					continue
				}

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

					swagType := parseGOTypeToSwaggerType(value.Kind())
					var optionalFormat string
					if swagType == "array" && isByteArray(field) {
						swagType = "string"
						optionalFormat = "binary"
					} else if swagType == "object" && isTime(field) {
						swagType = "string"
						optionalFormat = "date-time"
					}

					parameters = append(parameters, Parameter{
						Name:        fName,
						In:          string(qr.Type),
						Description: field.Tag.Get("description"),
						Required:    field.Tag.Get("required") == "true",
						Schema:      Schema{Type: swagType, Format: optionalFormat},
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

					bodyIsArray := isArray(br.Data)

					if bodyIsArray {
						for _, contentType := range br.ContentType {
							body.Content[contentType] = Content{
								Schema: Schema{
									Type: "array",
									Items: &Schema{
										Ref: fmt.Sprintf("#/components/schemas/%s", splitTypeName[len(splitTypeName)-1]),
									},
								},
							}
						}
					} else {
						for _, contentType := range br.ContentType {
							body.Content[contentType] = Content{
								Schema: Schema{
									Ref: fmt.Sprintf("#/components/schemas/%s", splitTypeName[len(splitTypeName)-1]),
								},
							}
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

					responseIsArray := isArray(res.Data)

					if responseIsArray {
						for _, contentType := range res.ContentType {
							content[contentType] = Content{
								Schema: Schema{
									Type: "array",
									Items: &Schema{
										Ref: fmt.Sprintf("#/components/schemas/%s", splitTypeName[len(splitTypeName)-1]),
									},
								},
							}
						}
					} else {
						for _, contentType := range res.ContentType {
							content[contentType] = Content{
								Schema: Schema{
									Ref: fmt.Sprintf("#/components/schemas/%s", splitTypeName[len(splitTypeName)-1]),
								},
							}
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

			securityMemberships := []map[string][]string{}

			if rd.AuthenticationConfiguration != nil {
				if rd.AuthenticationConfiguration.BasicAuth != nil {
					if rd.AuthenticationConfiguration.BasicAuth.Name == "" {
						rd.AuthenticationConfiguration.BasicAuth.Name = "basic"
					}
					securityMemberships = append(securityMemberships, map[string][]string{
						rd.AuthenticationConfiguration.BasicAuth.Name: {},
					})
				}
				if rd.AuthenticationConfiguration.BearerAuth != nil {
					if rd.AuthenticationConfiguration.BearerAuth.Name == "" {
						rd.AuthenticationConfiguration.BearerAuth.Name = "bearer"
					}
					securityMemberships = append(securityMemberships, map[string][]string{
						rd.AuthenticationConfiguration.BearerAuth.Name: {},
					})
				}
				if rd.AuthenticationConfiguration.ApiKeyAuth != nil {
					if rd.AuthenticationConfiguration.ApiKeyAuth.Name == "" {
						rd.AuthenticationConfiguration.ApiKeyAuth.Name = "apiKey"
					}
					securityMemberships = append(securityMemberships, map[string][]string{
						rd.AuthenticationConfiguration.ApiKeyAuth.Name: {},
					})
				}
				if rd.AuthenticationConfiguration.OpenIdAuth != nil {
					if rd.AuthenticationConfiguration.OpenIdAuth.Name == "" {
						rd.AuthenticationConfiguration.OpenIdAuth.Name = "openId"
					}
					securityMemberships = append(securityMemberships, map[string][]string{
						rd.AuthenticationConfiguration.OpenIdAuth.Name: {},
					})
				}
				if rd.AuthenticationConfiguration.Oauth2Auth != nil {
					if rd.AuthenticationConfiguration.Oauth2Auth.Name == "" {
						rd.AuthenticationConfiguration.Oauth2Auth.Name = "oauth2"
					}
					securityMemberships = append(securityMemberships, map[string][]string{
						rd.AuthenticationConfiguration.Oauth2Auth.Name: rd.OauthScopes,
					})
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
				Security:    securityMemberships,
			}
		}

	}

	return paths, nil
}

func (c *SwaggoMux) getSecuritySchemas() map[string]SecurityScheme {
	allAuthenticationConfigurations := ext.Where(ext.SliceMap(ext.FlattenMap(c.routes, func(route Route) []RequestDetails {
		return route.RequestDetails
	}), func(requestDetails RequestDetails) *AuthenticationConfiguration {
		return requestDetails.AuthenticationConfiguration
	}), func(authConfigPtr *AuthenticationConfiguration) bool {
		return authConfigPtr != nil
	})

	uniqueBasicAuthConfigurations := ext.DistinctBy(ext.Where(ext.SliceMap(allAuthenticationConfigurations, func(authConfigPtr *AuthenticationConfiguration) *BasicAuth {
		return authConfigPtr.BasicAuth
	}), func(basicAuthPtr *BasicAuth) bool {
		return basicAuthPtr != nil
	}), func(basicAuthPtr *BasicAuth) string {
		if basicAuthPtr.Name == "" {
			basicAuthPtr.Name = "basic"
		}
		return basicAuthPtr.Name
	})

	uniqueBearerConfigurations := ext.DistinctBy(ext.Where(ext.SliceMap(allAuthenticationConfigurations, func(authConfigPtr *AuthenticationConfiguration) *BearerAuth {
		return authConfigPtr.BearerAuth
	}), func(bearerAuthPtr *BearerAuth) bool {
		return bearerAuthPtr != nil
	}), func(bearerAuthPtr *BearerAuth) string {
		if bearerAuthPtr.Name == "" {
			bearerAuthPtr.Name = "bearer"
		}
		return bearerAuthPtr.Name
	})

	uniqueApiKeyConfigurations := ext.DistinctBy(ext.Where(ext.SliceMap(allAuthenticationConfigurations, func(authConfigPtr *AuthenticationConfiguration) *ApiKeyAuth {
		return authConfigPtr.ApiKeyAuth
	}), func(apiKeyAuthPtr *ApiKeyAuth) bool {
		return apiKeyAuthPtr != nil
	}), func(apiKeyAuthPtr *ApiKeyAuth) string {
		if apiKeyAuthPtr.Name == "" {
			apiKeyAuthPtr.Name = "apiKey"
		}
		return apiKeyAuthPtr.Name
	})

	uniqueOpenIdConfigurations := ext.DistinctBy(ext.Where(ext.SliceMap(allAuthenticationConfigurations, func(authConfigPtr *AuthenticationConfiguration) *OpenIdAuth {
		return authConfigPtr.OpenIdAuth
	}), func(openIdAuthPtr *OpenIdAuth) bool {
		return openIdAuthPtr != nil
	}), func(openIdAuthPtr *OpenIdAuth) string {
		if openIdAuthPtr.Name == "" {
			openIdAuthPtr.Name = "openId"
		}
		return openIdAuthPtr.Name
	})

	uniqueOauth2Configurations := ext.DistinctBy(ext.Where(ext.SliceMap(allAuthenticationConfigurations, func(authConfigPtr *AuthenticationConfiguration) *Oauth2Auth {
		return authConfigPtr.Oauth2Auth
	}), func(oauth2AuthPtr *Oauth2Auth) bool {
		return oauth2AuthPtr != nil
	}), func(oauth2AuthPtr *Oauth2Auth) string {
		if oauth2AuthPtr.Name == "" {
			oauth2AuthPtr.Name = "oauth2"
		}
		return oauth2AuthPtr.Name
	})

	securitySchemes := map[string]SecurityScheme{}

	for _, basicAuth := range uniqueBasicAuthConfigurations {
		if basicAuth.Name == "" {
			basicAuth.Name = "basic"
		}
		securitySchemes[basicAuth.Name] = SecurityScheme{
			Type:   "http",
			Scheme: "basic",
		}
	}

	for _, bearerAuth := range uniqueBearerConfigurations {
		if bearerAuth.Name == "" {
			bearerAuth.Name = "bearer"
		}
		securitySchemes[bearerAuth.Name] = SecurityScheme{
			Type:   "http",
			Scheme: "bearer",
		}
	}

	for _, apiKeyAuth := range uniqueApiKeyConfigurations {
		if apiKeyAuth.Name == "" {
			apiKeyAuth.Name = "apiKey"
		}
		securitySchemes[apiKeyAuth.Name] = SecurityScheme{
			Type: "apiKey",
			Name: apiKeyAuth.Name,
			In:   apiKeyAuth.In,
		}
	}

	for _, openIdAuth := range uniqueOpenIdConfigurations {
		if openIdAuth.Name == "" {
			openIdAuth.Name = "openId"
		}
		securitySchemes[openIdAuth.Name] = SecurityScheme{
			Type:             "openIdConnect",
			OpenIdConnectUrl: openIdAuth.OpenIdConnectUrl,
		}
	}

	for _, oauth2Auth := range uniqueOauth2Configurations {

		if oauth2Auth.Name == "" {
			oauth2Auth.Name = "oauth2"
		}

		var flows Flows

		if oauth2Auth.Flows.Implicit != nil {
			flows.Implicit = &Flow{
				AuthorizationURL: oauth2Auth.Flows.Implicit.AuthorizationUrl,
				Scopes:           oauth2Auth.Flows.Implicit.Scopes,
			}
		}

		if oauth2Auth.Flows.Password != nil {
			flows.Password = &Flow{
				AuthorizationURL: oauth2Auth.Flows.Password.AuthorizationUrl,
				TokenURL:         oauth2Auth.Flows.Password.TokenUrl,
				Scopes:           oauth2Auth.Flows.Password.Scopes,
			}
		}

		if oauth2Auth.Flows.ClientCredentials != nil {
			flows.ClientCredentials = &Flow{
				AuthorizationURL: oauth2Auth.Flows.ClientCredentials.AuthorizationUrl,
				TokenURL:         oauth2Auth.Flows.ClientCredentials.TokenUrl,
				Scopes:           oauth2Auth.Flows.ClientCredentials.Scopes,
			}
		}

		if oauth2Auth.Flows.AuthorizationCode != nil {
			flows.AuthorizationCode = &Flow{
				AuthorizationURL: oauth2Auth.Flows.AuthorizationCode.AuthorizationUrl,
				TokenURL:         oauth2Auth.Flows.AuthorizationCode.TokenUrl,
				Scopes:           oauth2Auth.Flows.AuthorizationCode.Scopes,
			}
		}

		securitySchemes[oauth2Auth.Name] = SecurityScheme{
			Type:  "oauth2",
			Flows: &flows,
		}
	}

	return securitySchemes
}
