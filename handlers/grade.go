package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/database"
)

func GetStudentGradesHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Временно используем student_id = 1
	studentID := 1

	rows, err := database.DB.Query(`
        SELECT к.Название, COALESCE(у.Оценка, 'не оценено'), to_char(у.Дата_оценки, 'DD.MM.YYYY')
        FROM Успеваемость у
        JOIN Курс к ON у.id_курса = к.id_курса
        WHERE у.id_студента = $1
    `, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var grades []map[string]interface{}
	for rows.Next() {
		var course, grade, date string
		err := rows.Scan(&course, &grade, &date)
		if err != nil {
			continue
		}
		grades = append(grades, map[string]interface{}{
			"course": course,
			"grade":  grade,
			"date":   date,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

func UpdateGradeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии и роли
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Преподаватель" {
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

	_, err := database.DB.Exec(`
        UPDATE Успеваемость 
        SET Оценка = $1, Дата_оценки = CURRENT_DATE
        WHERE id_студента = $2 AND id_курса = $3
    `, grade, studentID, courseID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func GetCourseStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии и роли
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) != "Преподаватель" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		http.Error(w, "course_id required", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
        SELECT с.id_студента, ф.Фамилия, ф.Имя, ф.Отчество, COALESCE(у.Оценка, '')
        FROM Успеваемость у
        JOIN Студент с ON у.id_студента = с.id_студента
        JOIN Физическое_лицо ф ON с.id_лица = ф.id_лица
        WHERE у.id_курса = $1
    `, courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []map[string]interface{}
	for rows.Next() {
		var id int
		var lastName, firstName, middleName, grade string
		err := rows.Scan(&id, &lastName, &firstName, &middleName, &grade)
		if err != nil {
			continue
		}
		students = append(students, map[string]interface{}{
			"id":        id,
			"full_name": lastName + " " + firstName + " " + middleName,
			"grade":     grade,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func GetAllGradesHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`
        SELECT 
            ф.Фамилия || ' ' || ф.Имя || ' ' || COALESCE(ф.Отчество, '') as student,
            к.Название as course,
            COALESCE(у.Оценка, 'не оценено') as grade,
            COALESCE(to_char(у.Дата_оценки, 'DD.MM.YYYY'), '-') as date
        FROM Успеваемость у
        JOIN Студент с ON у.id_студента = с.id_студента
        JOIN Физическое_лицо ф ON с.id_лица = ф.id_лица
        JOIN Курс к ON у.id_курса = к.id_курса
        ORDER BY с.id_студента, к.id_курса
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var grades []map[string]interface{}
	for rows.Next() {
		var student, course, grade, date string
		err := rows.Scan(&student, &course, &grade, &date)
		if err != nil {
			continue
		}
		grades = append(grades, map[string]interface{}{
			"student": student,
			"course":  course,
			"grade":   grade,
			"date":    date,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}
