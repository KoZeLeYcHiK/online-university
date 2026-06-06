package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/constants"
	"online-university/services"
	"strconv"
)

type CourseHandler struct {
	service     *services.CourseService
	authHandler *AuthHandler
}

func NewCourseHandler(service *services.CourseService, auth *AuthHandler) *CourseHandler {
	return &CourseHandler{service: service, authHandler: auth}
}

func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *CourseHandler) GetTeacherCourses(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != constants.RoleTeacher {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	courses, err := h.service.GetTeacherCourses(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *CourseHandler) AddCourse(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != constants.RoleAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	credits, _ := strconv.Atoi(r.FormValue("credits"))
	maxStudents, _ := strconv.Atoi(r.FormValue("max_students"))
	teacherID, _ := strconv.Atoi(r.FormValue("teacher_id"))

	err := h.service.Create(title, description, credits, maxStudents, teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != constants.RoleAdmin {
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

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != constants.RoleAdmin {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	title := r.FormValue("title")
	description := r.FormValue("description")
	credits, _ := strconv.Atoi(r.FormValue("credits"))
	maxStudents, _ := strconv.Atoi(r.FormValue("max_students"))
	teacherID, _ := strconv.Atoi(r.FormValue("teacher_id"))

	err := h.service.Update(id, title, description, credits, maxStudents, teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
