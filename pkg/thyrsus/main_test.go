package thyrsus_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/cedrus-and-thuja/thyrus-go/pkg/thyrsus"
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
