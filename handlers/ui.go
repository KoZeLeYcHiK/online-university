package handlers

import (
	"html/template"
	"net/http"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if !CheckRole(r, "Администратор") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, _ := template.ParseFiles("templates/admin.html")
	tmpl.Execute(w, map[string]string{"SessionID": r.URL.Query().Get("session_id")})
}

func TeacherPage(w http.ResponseWriter, r *http.Request) {
	if !CheckRole(r, "Преподаватель") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, _ := template.ParseFiles("templates/teacher.html")
	tmpl.Execute(w, map[string]string{"SessionID": r.URL.Query().Get("session_id")})
}

func StudentPage(w http.ResponseWriter, r *http.Request) {
	if !CheckRole(r, "Студент") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, _ := template.ParseFiles("templates/student.html")
	tmpl.Execute(w, map[string]string{"SessionID": r.URL.Query().Get("session_id")})
}
