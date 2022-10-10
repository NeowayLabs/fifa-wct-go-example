//go:build unit

package rest_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/application"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/domain"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/infrastructure/serve/rest"
)

func Test_handler_postTeam(t *testing.T) {
	t.Run("should post a team successfully", func(t *testing.T) {
		service := &application.TeamServiceMock{
			CreateFn: func(ctx context.Context, input *application.TeamInput) (*application.TeamOutput, error) {
				assert.NotNil(t, ctx)

				expected := application.NewTeamInput("brazil", "Brazil", "A")
				assert.EqualValues(t, expected, input)

				team, err := domain.NewTeam("brazil", "Brazil", "A")
				assert.NoError(t, err)

				return application.NewTeamOutput(team), nil
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)
		bodyBytes := `{"id":"brazil", "name":"Brazil", "group":"A"}`

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/teams", URL), strings.NewReader(bodyBytes))
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"id":"brazil", "name":"Brazil", "group":"A"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.CreateInvokedCount)
	})

	t.Run("should return handled error when invalid body", func(t *testing.T) {
		handler := rest.NewHandler(nil, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)
		bodyBytes := `{"id":"brazil"`

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/teams", URL), strings.NewReader(bodyBytes))
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"invalid body", "detail":"error to decode body: invalid body"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))
	})

	t.Run("should return handled error when service returns error", func(t *testing.T) {
		service := &application.TeamServiceMock{
			CreateFn: func(ctx context.Context, input *application.TeamInput) (*application.TeamOutput, error) {
				return nil, fmt.Errorf("error on create team: %w", domain.ErrInvalidArgument)
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)
		bodyBytes := `{"id":"brazil", "name":"Brazil", "group":"A"}`

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/teams", URL), strings.NewReader(bodyBytes))
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"invalid argument", "detail":"error on create team: invalid argument"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.CreateInvokedCount)
	})

	t.Run("should return handled error when panic", func(t *testing.T) {
		service := &application.TeamServiceMock{
			CreateFn: func(ctx context.Context, input *application.TeamInput) (*application.TeamOutput, error) {
				panic("service panic")
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)
		bodyBytes := `{"id":"brazil", "name":"Brazil", "group":"A"}`

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/teams", URL), strings.NewReader(bodyBytes))
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"internal server error", "detail":"the server encountered an unexpected condition that prevented it from fulfilling the request"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.CreateInvokedCount)
	})
}

func Test_handler_getTeamByID(t *testing.T) {
	t.Run("should get a team successfully", func(t *testing.T) {
		service := &application.TeamServiceMock{
			GetFn: func(ctx context.Context, ID string) (*application.TeamOutput, error) {
				assert.NotNil(t, ctx)
				assert.Equal(t, "brazil", ID)

				team, err := domain.NewTeam("brazil", "Brazil", "A")
				assert.NoError(t, err)

				return application.NewTeamOutput(team), nil
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/teams/brazil", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"id":"brazil", "name":"Brazil", "group":"A"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.GetInvokedCount)
	})

	t.Run("should return handled error when service returns error", func(t *testing.T) {
		service := &application.TeamServiceMock{
			GetFn: func(ctx context.Context, ID string) (*application.TeamOutput, error) {
				return nil, fmt.Errorf("error on get team by id: %w", domain.ErrInvalidArgument)
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/teams/id", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"invalid argument", "detail":"error on get team by id: invalid argument"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.GetInvokedCount)
	})

	t.Run("should return handled error when panic", func(t *testing.T) {
		service := &application.TeamServiceMock{
			GetFn: func(ctx context.Context, ID string) (*application.TeamOutput, error) {
				panic("service panic")
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/teams/id", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"internal server error", "detail":"the server encountered an unexpected condition that prevented it from fulfilling the request"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.GetInvokedCount)
	})
}

func Test_handler_deleteTeamByID(t *testing.T) {
	t.Run("should delete a team successfully", func(t *testing.T) {
		service := &application.TeamServiceMock{
			RemoveFn: func(ctx context.Context, ID string) error {
				return nil
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/teams/id", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Empty(t, string(bodyBytesResp))

		assert.Equal(t, 1, service.RemoveInvokedCount)
	})

	t.Run("should return handled error when service returns error", func(t *testing.T) {
		service := &application.TeamServiceMock{
			RemoveFn: func(ctx context.Context, ID string) error {
				return fmt.Errorf("error on delete team by id: %w", domain.ErrNotFound)
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/teams/id", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"not found", "detail":"error on delete team by id: not found"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.RemoveInvokedCount)
	})

	t.Run("should return handled error when panic", func(t *testing.T) {
		service := &application.TeamServiceMock{
			RemoveFn: func(ctx context.Context, ID string) error {
				panic("service panic")
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/teams/id", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"internal server error", "detail":"the server encountered an unexpected condition that prevented it from fulfilling the request"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.RemoveInvokedCount)
	})
}

func Test_handler_getTeams(t *testing.T) {
	t.Run("should get all teams successfully", func(t *testing.T) {
		service := &application.TeamServiceMock{
			GetAllFn: func(ctx context.Context) ([]*application.TeamOutput, error) {
				assert.NotNil(t, ctx)

				team, err := domain.NewTeam("brazil", "Brazil", "A")
				assert.NoError(t, err)

				return []*application.TeamOutput{
					application.NewTeamOutput(team),
				}, nil
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/teams", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{
			"teams": [{"id":"brazil", "name":"Brazil", "group":"A"}]
		}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.GetAllInvokedCount)
	})

	t.Run("should return handled error when service returns error", func(t *testing.T) {
		service := &application.TeamServiceMock{
			GetAllFn: func(ctx context.Context) ([]*application.TeamOutput, error) {
				return nil, fmt.Errorf("error on get all teams: %w", domain.ErrInvalidArgument)
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/teams", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"invalid argument", "detail":"error on get all teams: invalid argument"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.GetAllInvokedCount)
	})

	t.Run("should return handled error when panic", func(t *testing.T) {
		service := &application.TeamServiceMock{
			GetAllFn: func(ctx context.Context) ([]*application.TeamOutput, error) {
				panic("service panic")
			},
		}

		handler := rest.NewHandler(service, testLog)
		server := httptest.NewServer(handler)
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/teams", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `{"title":"internal server error", "detail":"the server encountered an unexpected condition that prevented it from fulfilling the request"}`
		assert.JSONEq(t, expectedBody, string(bodyBytesResp))

		assert.Equal(t, 1, service.GetAllInvokedCount)
	})
}
