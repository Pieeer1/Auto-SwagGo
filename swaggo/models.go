package swaggo

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestDataSource string

const (
	FormDataSource         RequestDataSource = "formData"
	QuerySource            RequestDataSource = "query"
	PathSource             RequestDataSource = "path"
	BodySource             RequestDataSource = "body"
	MultiformContentSource RequestDataSource = "multipart/form-data"
	HeaderSource           RequestDataSource = "header"
)

type Route struct {
	Path           string
	Handler        http.Handler
	Prefix         string
	Version        string
	RequestDetails []RequestDetails
}

func (r *Route) GetPathWithoutPrefixAndVersion() string {
	if r.Version == "" {
		return strings.Split(strings.Split(r.Path, fmt.Sprintf("%s/", r.Prefix))[1], "/")[0]
	}
	return strings.Split(strings.Split(r.Path, fmt.Sprintf("%s/%s/", r.Prefix, r.Version))[1], "/")[0]
}

type SwaggerInfo struct {
	Title                   string
	Description             string
	TermsOfServiceURL       string
	ContactEmail            string
	LicenseName             string
	LicenseURL              string
	Version                 string
	ExternalDocsDescription string
	ExternalDocsURL         string
	Servers                 []string
}

type RequestDetails struct {
	Method      string
	Summary     string
	Description string
	Requests    []RequestData
	Responses   []ResponseData
}

type RequestData struct {
	Type        RequestDataSource
	Description string
	Required    bool
	ContentType []string
	Data        any
}
type ResponseData struct {
	Code        int
	Data        any
	ContentType []string
	Headers     map[string]any
}
