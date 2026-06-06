package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/services"
)

type GradeHandler struct {
	service     *services.GradeService
	authHandler *AuthHandler
}

func NewGradeHandler(service *services.GradeService, auth *AuthHandler) *GradeHandler {
	return &GradeHandler{service: service, authHandler: auth}
}

func (h *GradeHandler) GetStudentGrades(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Студент" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	grades, err := h.service.GetStudentGrades(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

func (h *GradeHandler) GetAllGrades(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	grades, err := h.service.GetAllGrades()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

func (h *GradeHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Преподаватель" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentID := r.FormValue("student_id")
	courseID := r.FormValue("course_id")
	grade := r.FormValue("grade")

	err := h.service.UpdateGrade(studentID, courseID, grade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *GradeHandler) GetCourseStudents(w http.ResponseWriter, r *http.Request) {
	if h.authHandler.GetRole(r) != "Преподаватель" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		http.Error(w, "course_id required", http.StatusBadRequest)
		return
	}
	students, err := h.service.GetCourseStudents(courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
