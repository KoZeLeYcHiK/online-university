package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/database"
	"online-university/models"
)

func GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`
        SELECT к.id_курса, к.Название, к.Описание, к.Кредиты_ECTS, к.Макс_студентов,
               ф.Фамилия, ф.Имя, ф.Отчество
        FROM Курс к
        JOIN Преподаватель п ON к.id_преподавателя = п.id_преподавателя
        JOIN Физическое_лицо ф ON п.id_лица = ф.id_лица
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		var lastName, firstName, middleName string
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.CreditsECTS, &c.MaxStudents,
			&lastName, &firstName, &middleName)
		if err != nil {
			continue
		}
		c.TeacherName = lastName + " " + firstName + " " + middleName
		courses = append(courses, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func AddCourseHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	credits := r.FormValue("credits")
	maxStudents := r.FormValue("max_students")
	teacherID := r.FormValue("teacher_id")

	_, err := database.DB.Exec(`
        INSERT INTO Курс (Название, Описание, Кредиты_ECTS, Макс_студентов, id_преподавателя)
        VALUES ($1, $2, $3, $4, $5)
    `, title, description, credits, maxStudents, teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	courseID := r.URL.Query().Get("id")
	if courseID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("DELETE FROM Курс WHERE id_курса = $1", courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
