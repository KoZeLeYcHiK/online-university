package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/services"
	"strconv"
)

type StudentHandler struct {
	service     *services.StudentService
	authHandler *AuthHandler
}

func NewStudentHandler(service *services.StudentService, auth *AuthHandler) *StudentHandler {
	return &StudentHandler{service: service, authHandler: auth}
}

func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	students, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *StudentHandler) AddStudent(w http.ResponseWriter, r *http.Request) {
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
		r.FormValue("record_book"),
		r.FormValue("enrollment_year"),
		r.FormValue("status"),
		r.FormValue("direction_id"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
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

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
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
		r.FormValue("record_book"),
		r.FormValue("enrollment_year"),
		r.FormValue("status"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
