package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/application"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/domain"
)

type handler struct {
	teamService application.TeamService
	log         *log.Logger
}

func NewHandler(teamService application.TeamService, log *log.Logger) http.Handler {
	h := &handler{
		teamService: teamService,
		log:         log,
	}

	router := httprouter.New()

	// Default
	router.GET("/", h.info)

	// Teams
	router.POST("/teams", h.withRecover(h.postTeam))
	router.GET("/teams/:id", h.withRecover(h.getTeamByID))
	router.DELETE("/teams/:id", h.withRecover(h.deleteTeamByID))
	router.GET("/teams", h.withRecover(h.getTeams))

	return router
}

func (h *handler) withRecover(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if e := recover(); e != nil {
				h.writeJSONErr(w, fmt.Errorf("error recovered: %s", e)) //nolint:goerr113
			}
		}()
		next(w, r, ps)
	}
}

var errorToStatusCode = map[error]int{
	ErrInvalidBody:            http.StatusBadRequest,
	domain.ErrDuplicateKey:    http.StatusBadRequest,
	domain.ErrInvalidArgument: http.StatusBadRequest,
	domain.ErrNotFound:        http.StatusNotFound,
}

func (h *handler) writeJSONErr(rw http.ResponseWriter, err error) {
	for e, sc := range errorToStatusCode {
		if errors.Is(err, e) {
			h.writeProblemJSON(rw, sc, err)
			return
		}
	}

	h.log.Printf(err.Error())
	h.writeProblemJSON(rw, http.StatusInternalServerError, err)
}

type problemJSON struct {
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func (h *handler) writeProblemJSON(rw http.ResponseWriter, statusCode int, err error) {
	body := problemJSON{
		Title:  "internal server error",
		Detail: "the server encountered an unexpected condition that prevented it from fulfilling the request",
	}

	if cause := errors.Unwrap(err); cause != nil {
		body.Title = cause.Error()
		body.Detail = err.Error()
	}

	rw.Header().Set("Content-Type", "application/problem+json; charset=UTF-8")
	rw.WriteHeader(statusCode)

	if err := json.NewEncoder(rw).Encode(body); err != nil {
		h.log.Printf("error on encode response body: %v", err)
	}
}

func (h *handler) writeJSON(rw http.ResponseWriter, statusCode int, body interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(rw).Encode(body); err != nil {
			h.log.Printf("error on encode response body: %v", err)
		}
	}
}
