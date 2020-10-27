package client_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tufin/sabik/api"
	"github.com/tufin/sabik/client"
	"github.com/tufin/sabik/common/env"
)

func TestNewMiddleware(t *testing.T) {

	const domain, project, serviceName, headerKey, headerValue = "generic-bank", "retail", "customer", "hola", "mundo"
	const requestBody, responseBody = `{"cherry": true}`, `{"type": "coral"}`
	parsedURL, err := url.Parse("/ping?hello=world&hola=mundo")
	require.NoError(t, err)
	persister := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string][]*api.HTTPLog
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		httpLog := payload["logs"][0]
		require.Len(t, payload["logs"], 1)
		require.Equal(t, domain, httpLog.Domain)
		require.Equal(t, project, httpLog.Project)
		require.Equal(t, serviceName, httpLog.Destination)
		require.Equal(t, requestBody, httpLog.RequestBody)
		require.Equal(t, responseBody, httpLog.ResponseBody)
		requestHeaders := httpLog.GetRequestHeaders()
		require.Equal(t, "application/json", requestHeaders["Content-Type"][0])
		require.Equal(t, "me", requestHeaders["Test"][0])
		responseHeaders := httpLog.GetResponseHeaders()
		require.Equal(t, headerValue, responseHeaders[headerKey])
		w.WriteHeader(http.StatusCreated)
	}))
	require.NoError(t, os.Setenv(env.KeyDomain, domain))
	require.NoError(t, os.Setenv(env.KeyProject, project))
	require.NoError(t, os.Setenv(client.EnvKeyServiceName, serviceName))
	require.NoError(t, os.Setenv(client.EnvKeyTufinURL, persister.URL))
	request, err := http.NewRequest(http.MethodGet, parsedURL.String(), bytes.NewReader([]byte(requestBody)))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("test", "me")
	client.NewMiddleware().Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add(headerKey, headerValue)
		if _, err := io.WriteString(w, responseBody); err != nil {
			log.Errorf("failed to write response body '%s' with '%v'", requestBody, err)
		}
	})).ServeHTTP(httptest.NewRecorder(), request)
}
