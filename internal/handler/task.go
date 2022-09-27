package handler

import (
	"encoding/json"
	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
	"net/http"
)

type task struct {
	srv domain.TaskService
}

func NewTask(srv domain.TaskService) *task {
	return &task{srv: srv}
}

func (h *task) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.Task

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
