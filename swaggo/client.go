package swaggo

import (
	"auto-swaggo/internal/ext"
	"encoding/json"
	"fmt"
	"net/http"
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
		client.HandleFunc("/swagger/index.html", client.swagger, []string{"GET"}, "")
		client.HandleFunc("/openapi.json", client.swaggerJson, []string{"GET"}, "")
	}

	for _, version := range versions {
		client.HandleFunc("/swagger/index.html", client.swagger, []string{"GET"}, version)
		client.HandleFunc("/openapi.json", client.swaggerJson, []string{"GET"}, version)
	}

	return client
}

func (m *SwaggoMux) Handle(path string, handler http.Handler, methods []string, version string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var fullPath string

	if version == "" {
		fullPath = fmt.Sprintf("%s%s", m.prefix, path)
	} else {
		fullPath = fmt.Sprintf("%s/%s%s", m.prefix, version, path)
	}

	m.routes = append(m.routes, Route{Methods: methods, Path: fullPath, Handler: handler, Prefix: m.prefix, Version: version})
	m.mux.Handle(fullPath, m.defaultMiddleware(handler, methods))

}

func (m *SwaggoMux) defaultMiddleware(handler http.Handler, methods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ext.Contains(methods, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		handler.ServeHTTP(w, r)
	})
}

func (m *SwaggoMux) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request), methods []string, version string, requestsAndResponses ...RequestOrResponse) {
	m.Handle(path, http.HandlerFunc(handler), methods, version)
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

		for _, method := range route.Methods {
			paths[route.Path][strings.ToLower(method)] = Path{
				Tags:        []string{tagName},
				Summary:     fmt.Sprintf("%s %s", method, tagName),
				Description: fmt.Sprintf("%s Operation for %s", method, route.Path),
				OperationID: fmt.Sprintf("%s-%s", method, route.Path),
				Parameters:  []Parameter{},
				Responses:   map[string]Response{},
				Security:    []map[string][]string{},
			}
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
	}

	return doc, nil
}
