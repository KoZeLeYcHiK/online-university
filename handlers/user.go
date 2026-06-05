package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/database"
	"online-university/models"
)

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии и роли
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	rows, err := database.DB.Query(`
        SELECT с.id_студента, с.Номер_зачетки, с.Год_поступления, с.Статус,
               ф.Фамилия, ф.Имя, ф.Отчество
        FROM Студент с
        JOIN Физическое_лицо ф ON с.id_лица = ф.id_лица
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.RecordBookNumber, &s.EnrollmentYear, &s.Status,
			&s.LastName, &s.FirstName, &s.MiddleName)
		if err != nil {
			continue
		}
		students = append(students, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии и роли
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	rows, err := database.DB.Query(`
        SELECT ф.Фамилия, ф.Имя, ф.Отчество, п.Кафедра, кв.Название, уч.Название
        FROM Преподаватель п
        JOIN Физическое_лицо ф ON п.id_лица = ф.id_лица
        JOIN Квалификация кв ON п.id_квалификации = кв.id_квалификации
        JOIN Ученая_степень уч ON п.id_ученой_степени = уч.id_степени
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teachers []map[string]interface{}
	for rows.Next() {
		var lastName, firstName, middleName, department, qualification, degree string
		err := rows.Scan(&lastName, &firstName, &middleName, &department, &qualification, &degree)
		if err != nil {
			continue
		}
		teachers = append(teachers, map[string]interface{}{
			"full_name":     lastName + " " + firstName + " " + middleName,
			"department":    department,
			"qualification": qualification,
			"degree":        degree,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

func GetTeacherCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии и роли
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Преподаватель" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	teacherID := 1

	rows, err := database.DB.Query(`
        SELECT к.id_курса, к.Название, к.Описание, к.Кредиты_ECTS, к.Макс_студентов
        FROM Курс к
        WHERE к.id_преподавателя = $1
    `, teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []map[string]interface{}
	for rows.Next() {
		var id int
		var title, description string
		var credits, maxStudents int
		err := rows.Scan(&id, &title, &description, &credits, &maxStudents)
		if err != nil {
			continue
		}
		courses = append(courses, map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"credits_ects": credits,
			"max_students": maxStudents,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lastName := r.FormValue("last_name")
	firstName := r.FormValue("first_name")
	middleName := r.FormValue("middle_name")
	recordBook := r.FormValue("record_book")
	enrollmentYear := r.FormValue("enrollment_year")
	status := r.FormValue("status")
	directionID := r.FormValue("direction_id")

	// Сначала добавляем в Физическое_лицо
	var personID int
	err := database.DB.QueryRow(`
        INSERT INTO Физическое_лицо (Фамилия, Имя, Отчество, Дата_рождения, Электронная_почта, Телефон)
        VALUES ($1, $2, $3, '2000-01-01', 'temp@temp.ru', '79000000000')
        RETURNING id_лица
    `, lastName, firstName, middleName).Scan(&personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем в Студент
	_, err = database.DB.Exec(`
        INSERT INTO Студент (Номер_зачетки, Год_поступления, Статус, id_направления, id_лица)
        VALUES ($1, $2, $3, $4, $5)
    `, recordBook, enrollmentYear, status, directionID, personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	studentID := r.URL.Query().Get("id")
	if studentID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("DELETE FROM Студент WHERE id_студента = $1", studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lastName := r.FormValue("last_name")
	firstName := r.FormValue("first_name")
	middleName := r.FormValue("middle_name")
	department := r.FormValue("department")
	qualificationID := r.FormValue("qualification_id")
	degreeID := r.FormValue("degree_id")

	// Добавляем в Физическое_лицо
	var personID int
	err := database.DB.QueryRow(`
        INSERT INTO Физическое_лицо (Фамилия, Имя, Отчество, Дата_рождения, Электронная_почта, Телефон)
        VALUES ($1, $2, $3, '1980-01-01', 'temp@temp.ru', '79000000000')
        RETURNING id_лица
    `, lastName, firstName, middleName).Scan(&personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем в Преподаватель
	_, err = database.DB.Exec(`
        INSERT INTO Преподаватель (Кафедра, id_квалификации, id_ученой_степени, id_лица)
        VALUES ($1, $2, $3, $4)
    `, department, qualificationID, degreeID, personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Администратор" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	teacherID := r.URL.Query().Get("id")
	if teacherID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("DELETE FROM Преподаватель WHERE id_преподавателя = $1", teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
