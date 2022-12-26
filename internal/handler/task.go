package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PoorMercymain/REST-API-work-duration-counter/internal/domain"
	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/router"
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

func (h *task) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *task) Update(w http.ResponseWriter, r *http.Request) {
	var data domain.Task

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.srv.Update(r.Context(), data.Id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *task) ListWorksOfTask(w http.ResponseWriter, r *http.Request) {
	id, err := router.Params(r).Uint32("id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	works, err := h.srv.ListWorksOfTask(r.Context(), domain.Id(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = reply(w, works); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *task) CountDuration(w http.ResponseWriter, r *http.Request) {
	id, err := router.Params(r).Uint32("id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := h.srv.CountDuration(r.Context(), domain.Id(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = reply(w, result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *task) CreateTestTasks(w http.ResponseWriter, r *http.Request) {
	err := h.srv.CreateTestTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *task) CountAll(w http.ResponseWriter, r *http.Request) {
	duration, path, err := h.srv.CountAllDuration(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path = fmt.Sprintf("%d, %s", duration, path)

	if err = reply(w, path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//TODO: implement methods
