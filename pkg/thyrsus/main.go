package thyrsus

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var logger *log.Logger = log.New(os.Stderr, "", 0)

var EmptyHeaders map[string]string = make(map[string]string, 0)

type ExpectHttp struct {
	requests    []*MockRequest
	request_log []*http.Request
	Errors      []string
	testServer  *httptest.Server
}

type MockResponse interface {
	Status() int
	Body() []byte
	Headers() map[string]string
}

type BaseMockResponse struct {
	status  int
	body    []byte
	headers map[string]string
}

func (res BaseMockResponse) Status() int {
	return res.status
}

func (res BaseMockResponse) Body() []byte {
	return res.body
}

func (res BaseMockResponse) Headers() map[string]string {
	return res.headers
}

type MockRequest struct {
	Url      string
	Response MockResponse
}

func JSONMockResponse(status int, body interface{}, headers map[string]string) MockResponse {
	data, err := json.Marshal(body)
	if err != nil {
		data = []byte(err.Error())
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["ContentType"] = "application/json"
	return MockResponse(BaseMockResponse{
		status:  status,
		body:    data,
		headers: headers,
	})
}

func NewExpectHttp() *ExpectHttp {
	return &ExpectHttp{
		requests:    make([]*MockRequest, 0),
		request_log: make([]*http.Request, 0),
	}
}

func (exp *ExpectHttp) ExpectRequest(url string, response MockResponse) {
	exp.requests = append(exp.requests, &MockRequest{Url: url, Response: response})
}

func (exp *ExpectHttp) Start() string {
	exp.testServer = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exp.request_log = append(exp.request_log, r)
		if len(exp.requests) == 0 {
			exp.Errors = append(exp.Errors, fmt.Sprintf("no more requests expected, but got %s", r.RequestURI))
			w.WriteHeader(500)
		} else {
			expected := exp.requests[0]
			if expected.Url != r.RequestURI {
				exp.Errors = append(exp.Errors, fmt.Sprintf("expected %s, but got %s", expected.Url, r.RequestURI))
				w.WriteHeader(500)
			} else {
				exp.requests = exp.requests[1:]
				logger.Printf("request received: %s", r.RequestURI)
				for k, v := range expected.Response.Headers() {
					logger.Printf("setting response header: %s", k)
					w.Header().Set(k, v)
				}
				w.WriteHeader(expected.Response.Status())
				w.Write(expected.Response.Body())
			}
		}
	}))
	exp.testServer.EnableHTTP2 = true
	exp.testServer.Start()
	return exp.testServer.URL
}

func (exp *ExpectHttp) Close() {
	exp.testServer.Close()
}

func (exp *ExpectHttp) ValidateExpectations(T *testing.T) {
	for _, err := range exp.Errors {
		T.Error(err)
	}
}
