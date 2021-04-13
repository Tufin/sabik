package api_test

import (
	"net/http"
	"net/url"
	"testing"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/require"
	"github.com/tufin/sabik/api"
)

func getLog(t *testing.T) *api.HTTPLog {
	time, err := civil.ParseDateTime("2006-01-02t15:04:05.999999999")
	require.NoError(t, err)

	queryString, err := url.ParseQuery("a=b")
	require.NoError(t, err)

	requestHeaders := http.Header{}
	requestHeaders.Add("a", "b")

	responseHeaders := http.Header{}
	responseHeaders.Add("a", "b")

	cookies := []*http.Cookie{
		{
			Name:     "a",
			Value:    "b",
			Unparsed: []string{"a"},
		},
	}

	return &api.HTTPLog{
		Domain:          "a",
		Project:         "a",
		Time:            time,
		Scheme:          "a",
		Host:            "a",
		Path:            "a",
		QueryString:     queryString,
		Method:          "a",
		RequestBody:     "a",
		RequestHeaders:  requestHeaders,
		Cookies:         cookies,
		StatusCode:      1,
		ResponseBody:    "a",
		ResponseHeaders: responseHeaders,
		Service:         "a",
		Protocol:        "a",
		Connection:      "a",
		Origin:          "a",
	}
}

func TestClone_RequestHeaders(t *testing.T) {

	log := getLog(t)
	copy := log
	clone := log.Clone()

	delete(log.RequestHeaders, "A")

	require.NotEqual(t, getLog(t).RequestHeaders, copy.RequestHeaders)
	require.Equal(t, getLog(t), clone)
}

func TestClone_ResponseHeaders(t *testing.T) {

	log := getLog(t)
	copy := log
	clone := log.Clone()

	log.ResponseHeaders.Add("x", "y")

	require.NotEqual(t, getLog(t), copy)
	require.Equal(t, getLog(t), clone)
}

func TestClone_QueryString(t *testing.T) {

	log := getLog(t)
	copy := log
	clone := log.Clone()

	log.QueryString.Del("a")

	require.NotEqual(t, getLog(t), copy)
	require.Equal(t, getLog(t), clone)
}

func TestClone_Cookies(t *testing.T) {

	log := getLog(t)
	copy := log
	clone := log.Clone()

	log.Cookies[0].Unparsed[0] = "x"

	newLog := getLog(t)
	require.NotEqual(t, newLog.Cookies[0].Unparsed, copy.Cookies[0].Unparsed)
	require.Equal(t, newLog.Cookies[0].Unparsed, clone.Cookies[0].Unparsed)
}
