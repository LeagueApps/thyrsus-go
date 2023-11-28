package thyrsus_test

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"testing"

	"github.com/cedrus-and-thuja/thyrsus-go/pkg/thyrsus"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Status string `json:"status"`
}

func TestBasicJSON(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	url := expect.Start()
	defer expect.Close()
	fmt.Printf("url: %s", url)
	expect.ExpectRequest("/api/v1/sunshine", thyrsus.JSONMockResponse(200, &TestData{Status: "ok"}, thyrsus.EmptyHeaders))
	res, err := http.DefaultClient.Get(fmt.Sprintf("%s/api/v1/sunshine", url))
	assert.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "{\"status\":\"ok\"}", string(body))
}

func TestBasicJSONTooManyRequests(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	url := expect.Start()
	defer expect.Close()
	fmt.Printf("url: %s", url)
	expect.ExpectRequest("/api/v1/sunshine", thyrsus.JSONMockResponse(200, &TestData{Status: "ok"}, thyrsus.EmptyHeaders))
	res, err := http.DefaultClient.Get(fmt.Sprintf("%s/api/v1/sunshine", url))
	assert.Equal(t, 200, res.StatusCode, "http status code incorrect")
	assert.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "{\"status\":\"ok\"}", string(body))
	res, err = http.DefaultClient.Get(fmt.Sprintf("%s/api/v1/moonlight", url))
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode, "http status code incorrect")
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, "no more requests expected, but got /api/v1/moonlight", string(expect.Errors[0]))
}

func TestBasicJSONMismatchURL(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	url := expect.Start()
	defer expect.Close()
	fmt.Printf("url: %s", url)
	expect.ExpectRequest("/api/v1/sunshine", thyrsus.JSONMockResponse(200, &TestData{Status: "ok"}, thyrsus.EmptyHeaders))
	_, err := http.DefaultClient.Get(fmt.Sprintf("%s/api/v1/moonlight", url))
	assert.NoError(t, err)
	assert.Equal(t, "expected /api/v1/sunshine, but got /api/v1/moonlight", string(expect.Errors[0]))
}

func TestNoFail(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	expect.ValidateExpectations(t)
}

func TestRequestHeaders(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	url := expect.Start()
	defer expect.Close()
	headers := map[string]string{"Authorization": "bearer 8888"}
	expect.ExpectRequestWithHeaders("/api/flatulence/v1/fart", thyrsus.JSONMockResponse(200, &TestData{Status: "ok"}, thyrsus.EmptyHeaders), headers)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/flatulence/v1/fart", url), nil)
	for k, v := range headers {
		req.Header[k] = []string{v}
	}
	assert.NoError(t, err, "reguest not made")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	// assert.Equal(t, "expected /api/v1/sunshine, but got /api/v1/moonlight", string(expect.Errors[0]))

	expect.ValidateExpectations(t)
}

func TestRequestHeadersMissingHeader(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	url := expect.Start()
	defer expect.Close()
	headers := map[string]string{"Authorization": "bearer 8888"}
	expect.ExpectRequestWithHeaders("/api/flatulence/v1/fart", thyrsus.JSONMockResponse(200, &TestData{Status: "ok"}, thyrsus.EmptyHeaders), headers)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/flatulence/v1/fart", url), nil)
	assert.NoError(t, err, "reguest not made")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, "expected Authorization header to contain bearer 8888, it did not", string(expect.Errors[0]))
}

func TestFailJSONMarshal(t *testing.T) {
	expect := thyrsus.NewExpectHttp()
	url := expect.Start()
	defer expect.Close()
	expect.ExpectRequest("/api/v1/moonlight", thyrsus.JSONMockResponse(200, math.NaN(), thyrsus.EmptyHeaders))
	res, err := http.DefaultClient.Get(fmt.Sprintf("%s/api/v1/moonlight", url))
	assert.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "json: unsupported value: NaN", string(body))
}
