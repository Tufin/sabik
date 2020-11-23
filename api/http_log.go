package api

import (
	"net/http"
	"net/url"

	"cloud.google.com/go/civil"
)

const (
	OriginAWSAPIGateway   = "AWSAPIGateway"
	OriginIstioAccessLogs = "IstioAccessLogs"
	OriginIstioEnvoyLua   = "IstioEnvoyLua"
	OriginNginx           = "Nginx"
	OriginGoMiddleware    = "GoMiddleware"
)

type HTTPLog struct {
	Domain          string         `json:"domain"`
	Project         string         `json:"project"`
	Time            civil.DateTime `json:"time"`
	Scheme          string         `json:"scheme"`
	Host            string         `json:"host"` // host or host:port
	Path            string         `json:"path"`
	QueryString     url.Values     `json:"query_string"` // parsed (not encoded)
	Method          string         `json:"method"`
	RequestBody     string         `json:"request_body"`
	RequestHeaders  http.Header    `json:"request_headers"` // canonical format
	Cookies         []*http.Cookie `json:"cookies"`
	StatusCode      int            `json:"status_code"`
	ResponseBody    string         `json:"response_body"`
	ResponseHeaders http.Header    `json:"response_headers"` // canonical format
	Service         string         `json:"service"`
	Protocol        string         `json:"protocol"`
	Origin          string         `json:"origin"`
}
