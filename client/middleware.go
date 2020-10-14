package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/sabik/api"
	"github.com/tufin/sabik/common"
	"github.com/tufin/sabik/common/env"
)

type Middleware struct {
	domain      string
	project     string
	serviceName string
	tufinURL    string
}

const (
	EnvKeyServiceName = "TUFIN_SERVICE_NAME"
	EnvKeyTufinURL    = "TUFIN_URL"
)

func NewMiddleware() *Middleware {

	return &Middleware{
		domain:      env.GetEnvOrExit(env.KeyDomain),
		project:     env.GetEnvOrExit(env.KeyProject),
		serviceName: env.GetEnvOrExit(EnvKeyServiceName),
		tufinURL:    env.GetEnvWithDefault(EnvKeyTufinURL, "https://persister-xiixymmvca-ew.a.run.app"),
	}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := NewResponseRecorder(w)
		t := time.Now()
		next.ServeHTTP(recorder, r)
		report(m.domain, m.project, m.serviceName, m.tufinURL, recorder, r, t)
	})
}

func toHTTPLog(domain string, project string, serviceName string, response *ResponseRecorder, request *http.Request, t time.Time) (*api.HTTPLog, error) {

	requestHeader, err := common.StringHeader(request.Header)
	if err != nil {
		return nil, err
	}

	return &api.HTTPLog{
		Domain:          domain,
		Project:         project,
		Time:            civil.DateTimeOf(t),
		URI:             request.URL.Path,
		QueryString:     request.URL.RawQuery,
		Method:          request.Method,
		RequestBody:     getBody(request),
		RequestHeaders:  requestHeader,
		Source:          "",
		StatusCode:      response.StatusCode,
		ResponseBody:    response.Body.String(),
		ResponseHeaders: "",
		Destination:     serviceName,
		DestinationIP:   "",
		DestinationPort: 0,
		Protocol:        request.Proto,
		Client:          "",
		Origin:          api.OriginGoMiddleware,
	}, nil
}

func report(domain string, project string, serviceName string, tufinURL string, response *ResponseRecorder, request *http.Request, t time.Time) {

	httpLog, err := toHTTPLog(domain, project, serviceName, response, request, t)
	if err != nil {
		log.Errorf("failed to create HTTPLog with '%v'", err)
		return
	}

	body, err := json.Marshal(map[string][]*api.HTTPLog{"logs": {httpLog}})
	if err != nil {
		log.Errorf("failed to marshal HTTPLog with '%v'", err)
	} else {
		if response, err := http.Post(tufinURL, "application/json", bytes.NewReader(body)); err != nil {
			log.Errorf("failed to send HTTPLog '%s' to '%s' with '%v'", httpLog.URI, tufinURL, err)
		} else if response.StatusCode != http.StatusCreated {
			log.Errorf("failed to send HTTPLog '%s' with '%s'", httpLog.URI, response.Status)
		} else {
			log.Infof("sent HTTPLog '%s'", httpLog.URI)
		}
	}
}

type ResponseRecorder struct {
	writer     http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func NewResponseRecorder(writer http.ResponseWriter) *ResponseRecorder {

	return &ResponseRecorder{
		writer:     writer,
		Body:       new(bytes.Buffer),
		StatusCode: http.StatusOK,
	}
}

func (rw *ResponseRecorder) Header() http.Header {

	return rw.writer.Header()
}

func (rw *ResponseRecorder) Write(buf []byte) (int, error) {

	rw.Body.Write(buf)

	return rw.writer.Write(buf)
}

func (rw *ResponseRecorder) WriteHeader(statusCode int) {

	rw.StatusCode = statusCode
	rw.writer.WriteHeader(statusCode)
}

func getBody(r *http.Request) string {

	if r.Body == nil {
		return ""
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Errorf("failed to convert request body to string with '%v' url '%s'", err, r.URL.String())
		return ""
	}
	defer common.CloseWithErrLog(r.Body)

	return buf.String()
}
