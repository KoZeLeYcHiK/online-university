package handlers

import (
	"encoding/json"
	"net/http"
	"online-university/database"
)

func GetScheduleHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" || GetRole(r) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`
        SELECT к.Название, 
               ф.Фамилия || ' ' || ф.Имя || ' ' || COALESCE(ф.Отчество, '') as Преподаватель,
               to_char(р.Дата_время, 'DD.MM.YYYY HH24:MI') as Дата_время,
               COALESCE(р.Аудитория, 'не указана') as Аудитория,
               COALESCE(р.Ссылка_на_подключение, '-') as Ссылка
        FROM Расписание р
        JOIN Курс к ON р.id_курса = к.id_курса
        JOIN Преподаватель п ON р.id_преподавателя = п.id_преподавателя
        JOIN Физическое_лицо ф ON п.id_лица = ф.id_лица
        ORDER BY р.Дата_время
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var schedules []map[string]interface{}
	for rows.Next() {
		var course, teacher, datetime, room, link string
		err := rows.Scan(&course, &teacher, &datetime, &room, &link)
		if err != nil {
			continue
		}
		schedules = append(schedules, map[string]interface{}{
			"course":       course,
			"teacher":      teacher,
			"datetime":     datetime,
			"room":         room,
			"meeting_link": link,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}
