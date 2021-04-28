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
	BytesSent       int64          `json:"bytes_sent"`
	ResponseHeaders http.Header    `json:"response_headers"` // canonical format
	RequestTime     int64          `json:"request_time"`     // latency
	Service         string         `json:"service"`
	Protocol        string         `json:"protocol"`
	Connection      string         `json:"connection"`
	Origin          string         `json:"origin"`
}

var originResponses = map[string]bool{
	OriginAWSAPIGateway:   true,
	OriginIstioAccessLogs: true,
	OriginIstioEnvoyLua:   true,
	OriginNginx:           false,
	OriginGoMiddleware:    true,
}

func OriginHasResponse(origin string) bool {

	if r, ok := originResponses[origin]; ok {
		return r
	}

	return true
}

func (httpLog *HTTPLog) Clone() *HTTPLog {
	clone := *httpLog

	clone.QueryString = cloneURLValues(httpLog.QueryString)
	clone.RequestHeaders = httpLog.RequestHeaders.Clone()
	clone.Cookies = cloneCookies(httpLog.Cookies)
	clone.ResponseHeaders = httpLog.ResponseHeaders.Clone()

	return &clone
}

func cloneURLValues(v url.Values) url.Values {
	if v == nil {
		return nil
	}
	// http.Header and url.Values have the same representation, so temporarily
	// treat it like http.Header, which does have a clone:
	return url.Values(http.Header(v).Clone())
}

func cloneCookies(cookies []*http.Cookie) []*http.Cookie {
	if cookies == nil {
		return nil
	}

	result := make([]*http.Cookie, len(cookies))

	for i, cookie := range cookies {
		if cookie == nil {
			continue
		}

		newCookie := *cookie
		result[i] = &newCookie

		if cookie.Unparsed != nil {
			result[i].Unparsed = make([]string, len(cookie.Unparsed))
			copy(result[i].Unparsed, cookie.Unparsed)
		}
	}

	return result
}
