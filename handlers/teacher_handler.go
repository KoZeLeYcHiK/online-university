package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/services"
	"strconv"
)

type TeacherHandler struct {
	service     *services.TeacherService
	authHandler *AuthHandler
}

func NewTeacherHandler(service *services.TeacherService, auth *AuthHandler) *TeacherHandler {
	return &TeacherHandler{service: service, authHandler: auth}
}

func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	teachers, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

func (h *TeacherHandler) AddTeacher(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.service.Create(
		r.FormValue("last_name"),
		r.FormValue("first_name"),
		r.FormValue("middle_name"),
		r.FormValue("department"),
		r.FormValue("qualification_id"),
		r.FormValue("degree_id"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	err := h.service.Update(
		id,
		r.FormValue("department"),
		r.FormValue("qualification_id"),
		r.FormValue("degree_id"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
