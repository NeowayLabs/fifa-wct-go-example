package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/application"
)

type teamRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type teamResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type teamsResponse struct {
	Teams []teamResponse `json:"teams"`
}

func (h *handler) postTeam(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var teamRequest teamRequest
	if err := json.NewDecoder(r.Body).Decode(&teamRequest); err != nil {
		h.writeJSONErr(w, fmt.Errorf("error to decode body: %w", ErrInvalidBody))
		return
	}

	input := application.NewTeamInput(teamRequest.ID, teamRequest.Name, teamRequest.Group)

	output, err := h.teamService.Create(r.Context(), input)
	if err != nil {
		h.writeJSONErr(w, err)
		return
	}

	teamResponse := teamResponse{ID: output.ID, Name: output.Name, Group: output.Group}
	h.writeJSON(w, http.StatusCreated, teamResponse)
}

func (h *handler) getTeamByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	output, err := h.teamService.Get(r.Context(), ID)
	if err != nil {
		h.writeJSONErr(w, err)
		return
	}

	teamResponse := teamResponse{ID: output.ID, Name: output.Name, Group: output.Group}
	h.writeJSON(w, http.StatusOK, teamResponse)
}

func (h *handler) deleteTeamByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	err := h.teamService.Remove(r.Context(), ID)
	if err != nil {
		h.writeJSONErr(w, err)
		return
	}

	h.writeJSON(w, http.StatusNoContent, nil)
}

func (h *handler) getTeams(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	outputs, err := h.teamService.GetAll(r.Context())
	if err != nil {
		h.writeJSONErr(w, err)
		return
	}

	teams := make([]teamResponse, 0)

	for _, output := range outputs {
		team := teamResponse{ID: output.ID, Name: output.Name, Group: output.Group}
		teams = append(teams, team)
	}

	teamsResponse := teamsResponse{Teams: teams}
	h.writeJSON(w, http.StatusOK, teamsResponse)
}
