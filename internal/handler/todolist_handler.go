package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ankush263/todolist/internal/model"
	"github.com/Ankush263/todolist/internal/repository"
	"github.com/gorilla/mux"
)

type TodolistHandler struct {
	repo *repository.TodolistRepository
}

func NewTodolistHandler(repo *repository.TodolistRepository) *TodolistHandler {
	return &TodolistHandler{repo: repo}
}

func (h *TodolistHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.TodoList
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validation checks
	if strings.TrimSpace(t.Title) == "" {
		http.Error(w, "Bad Object: title", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TodolistHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 64, 10)

	todolist, err := h.repo.GetSingleTodolist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todolist)
}

func (h *TodolistHandler) Update(w http.ResponseWriter, r *http.Request) {
	var t model.TodoList
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 64, 10)

	resp, err := h.repo.Update(r.Context(), &t, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *TodolistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 64, 10)

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
