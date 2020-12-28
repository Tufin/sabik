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

const (
	EnvKeyServiceName = "TUFIN_SABIK_SERVICE_NAME"
	EnvKeyTufinURL    = "TUFIN_SABIK_URL"
	EnvKeyEnable      = "TUFIN_SABIK_ENABLE"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type EmptyMiddleware struct{}

func (sm *EmptyMiddleware) Handle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type SabikMiddleware struct {
	domain      string
	project     string
	serviceName string
	tufinURL    string
	enable      bool
}

func newSabikMiddleware() Middleware {

	return &SabikMiddleware{
		domain:      env.GetEnvOrExit(env.KeyDomain),
		project:     env.GetEnvOrExit(env.KeyProject),
		serviceName: env.GetEnvOrExit(EnvKeyServiceName),
		tufinURL:    env.GetEnvWithDefault(EnvKeyTufinURL, "https://persister-xiixymmvca-ew.a.run.app"),
		enable:      isEnable(),
	}
}

func (sm *SabikMiddleware) Handle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sm.enable {
			recorder := NewResponseRecorder(w)
			t := time.Now()
			next.ServeHTTP(recorder, r)
			report(sm.domain, sm.project, sm.serviceName, sm.tufinURL, recorder, r, t)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func toHTTPLog(domain string, project string, serviceName string, response *ResponseRecorder, request *http.Request, t time.Time) (*api.HTTPLog, error) {

	return &api.HTTPLog{
		Domain:          domain,
		Project:         project,
		Time:            civil.DateTimeOf(t),
		Scheme:          request.URL.Scheme,
		Host:            request.URL.Host,
		Path:            request.URL.Path,
		QueryString:     request.URL.Query(),
		Method:          request.Method,
		RequestBody:     getBody(request),
		RequestHeaders:  request.Header,
		StatusCode:      response.StatusCode,
		ResponseBody:    response.Body.String(),
		ResponseHeaders: response.Header(),
		Service:         serviceName,
		Protocol:        request.Proto,
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
			log.Errorf("failed to send HTTPLog '%s' to '%s' with '%v'", httpLog.Path, tufinURL, err)
		} else if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
			log.Errorf("failed to send HTTPLog '%s' with '%s'", httpLog.Path, response.Status)
		} else {
			log.Infof("sent HTTPLog '%+v'", *httpLog)
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
