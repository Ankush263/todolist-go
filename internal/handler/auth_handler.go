package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ankush263/todolist/internal/auth"
	"github.com/Ankush263/todolist/internal/model"
	"github.com/Ankush263/todolist/internal/repository"
)

type AuthHandler struct {
	users *repository.UserRepository
}

func NewAuthHandler(users *repository.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

type authrequest struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid Credential", http.StatusUnauthorized)
		return
	}

	if err = auth.CheckPassword(req.Password, user.Password); err != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Token Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req authrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Password Error", http.StatusInternalServerError)
		return
	}

	user := model.User{
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Password: hash,
	}

	if err := h.users.Create(r.Context(), &user); err != nil {
		http.Error(w, "Invalid Credentials", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
