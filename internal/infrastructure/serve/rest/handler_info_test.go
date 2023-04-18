//go:build unit

package rest_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/NeowayLabs/fifa-wct-go-example/internal/infrastructure/serve/rest"
	"github.com/stretchr/testify/assert"
)

func TestHandlerInfo(t *testing.T) {
	handler := rest.NewHandler(nil, testLog)
	server := httptest.NewServer(handler)
	defer server.Close()

	URL, _ := url.Parse(server.URL)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/", URL), nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	bodyBytesResp, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyBytesResp)

	expectedBody := `
		{
			"title": "FIFA World Cup Table (GO Example)",
			"description": "The responsibility of this project is to store the teams, games and results of the FIFA World Cup"
		}`

	assert.JSONEq(t, expectedBody, string(bodyBytesResp))
}
