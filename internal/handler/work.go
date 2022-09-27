package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
)

type work struct {
	srv domain.WorkService
}

func NewWork(srv domain.WorkService) *work {
	return &work{srv: srv}
}

func (h *work) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.Work

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.srv.Create(r.Context(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		ID domain.Id `json:"id"`
	}{ID: id}

	if err = reply(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//TODO: implement methods
