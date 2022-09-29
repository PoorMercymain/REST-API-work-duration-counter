package handler

import (
	"encoding/json"
	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/router"
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

func (h *work) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := router.Params(r).Uint32("id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.Delete(r.Context(), domain.Id(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *work) List(w http.ResponseWriter, r *http.Request) {
	id, err := router.Params(r).Uint32("id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	works, err := h.srv.List(r.Context(), domain.Id(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		WORKS []domain.Work `json:"works"`
	}{WORKS: append(works.Parental, works.Main)}

	if err = reply(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//TODO: implement methods
