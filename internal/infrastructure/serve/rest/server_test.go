//go:build unit

package rest_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/infrastructure/serve/rest"
)

func TestServerListenAndServeAndShutdown(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	server := rest.NewServer(handler, "", "8080", log)
	server.ListenAndServe()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = server.Shutdown(time.Second)
	assert.NoError(t, err)

	res, err = http.DefaultClient.Do(req)
	assert.Error(t, err)
}
