package swaggo

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestDataSource string
type RequestOrResponseType string

type RequestDetails struct {
	Request  RequestOrResponseType
	Response RequestOrResponseType
}

const (
	FormDataSource         RequestDataSource = "formData"
	QuerySource            RequestDataSource = "query"
	PathSource             RequestDataSource = "path"
	BodySource             RequestDataSource = "body"
	MultiformContentSource RequestDataSource = "multipart/form-data"
	HeaderSource           RequestDataSource = "header"
)

type Route struct {
	Path                 string
	Methods              []string
	Handler              http.Handler
	Prefix               string
	Version              string
	RequestsAndResponses []RequestOrResponse
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

type RequestOrResponse struct {
	Type              RequestOrResponseType
	RequestDataSource *RequestDataSource
	Data              any
}
