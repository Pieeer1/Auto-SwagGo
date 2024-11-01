package swaggo

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestDataSource string

const (
	QuerySource  RequestDataSource = "query"
	PathSource   RequestDataSource = "path"
	BodySource   RequestDataSource = "body"
	HeaderSource RequestDataSource = "header"
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
	Method                      string
	Summary                     string
	Description                 string
	AuthenticationConfiguration *AuthenticationConfiguration
	OauthScopes                 []string
	Requests                    []RequestData
	Responses                   []ResponseData
}

type AuthenticationConfiguration struct {
	BasicAuth  *BasicAuth
	BearerAuth *BearerAuth
	ApiKeyAuth *ApiKeyAuth
	OpenIdAuth *OpenIdAuth
	Oauth2Auth *Oauth2Auth
}

type BasicAuth struct {
	Name string
}

type BearerAuth struct {
	Name string
}

type ApiKeyAuth struct {
	In   string
	Name string
}

type OpenIdAuth struct {
	Name             string
	OpenIdConnectUrl string
}

type Oauth2Auth struct {
	Name        string
	Description string
	Flows       Oauth2Flows
}

type Oauth2Flows struct {
	Implicit          *Oauth2Flow
	Password          *Oauth2Flow
	ClientCredentials *Oauth2Flow
	AuthorizationCode *Oauth2Flow
}

type Oauth2Flow struct {
	AuthorizationUrl string
	TokenUrl         string
	RefreshUrl       string
	Scopes           map[string]string
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
