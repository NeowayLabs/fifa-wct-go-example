package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *handler) info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := map[string]string{
		"title":       "FIFA World Cup Table (GO Example)",
		"description": "The responsibility of this project is to store the teams, games and results of the FIFA World Cup",
	}

	h.writeJSON(w, http.StatusOK, info)
}
