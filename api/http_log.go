package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"cloud.google.com/go/civil"
)

const (
	OriginGoMiddleware = "GoMiddleware"
)

type HTTPLog struct {
	Domain          string         `bigquery:"domain" json:"domain"`
	Project         string         `bigquery:"project" json:"project"`
	Time            civil.DateTime `bigquery:"time" json:"time"`
	URI             string         `bigquery:"uri" json:"uri"`
	QueryString     string         `bigquery:"query_string" json:"query_string"`
	Method          string         `bigquery:"request_method" json:"request_method"`
	RequestBody     string         `bigquery:"request_body" json:"request_body"`
	RequestHeaders  string         `bigquery:"request_headers" json:"request_headers"`
	Source          string         `bigquery:"source" json:"source"`
	StatusCode      int            `bigquery:"status_code" json:"status_code"`
	ResponseBody    string         `bigquery:"response_body" json:"response_body"`
	ResponseHeaders string         `bigquery:"response_headers" json:"response_headers"`
	Destination     string         `bigquery:"destination" json:"destination"`
	DestinationIP   string         `bigquery:"destination_ip" json:"destination_ip"`
	DestinationPort int            `bigquery:"destination_port" json:"destination_port"`
	Protocol        string         `bigquery:"protocol" json:"protocol"`
	Client          string         `bigquery:"client" json:"client"`
	Origin          string         `bigquery:"origin" json:"origin"`
}

func (l *HTTPLog) GetQueryString() (url.Values, error) {

	if l.QueryString == "" {
		return url.Values{}, nil
	}

	return url.ParseQuery(l.QueryString)
}

func (l *HTTPLog) GetResponseHeaders() http.Header {

	return toHTTPHeader(l.ResponseHeaders)
}

func (l *HTTPLog) GetRequestHeaders() http.Header {

	return toHTTPHeader(l.RequestHeaders)
}

func toHTTPHeader(headers string) http.Header {

	var ret http.Header
	err := json.Unmarshal([]byte(headers), &ret)
	if err != nil {
		log.Fatalf("failed to unmarshal HTTPLog headers '%s' into 'http.Header' with '%v'", headers, err)
	}

	return ret
}
