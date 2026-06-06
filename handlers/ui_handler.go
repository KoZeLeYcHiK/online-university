package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

type UIHandler struct {
	db          *sql.DB
	authHandler *AuthHandler
}

func NewUIHandler(db *sql.DB, auth *AuthHandler) *UIHandler {
	return &UIHandler{db: db, authHandler: auth}
}

func (h *UIHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	if !h.authHandler.CheckRole(r, "Администратор") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, _ := template.ParseFiles("templates/admin.html")
	tmpl.Execute(w, map[string]string{"SessionID": r.URL.Query().Get("session_id")})
}

func (h *UIHandler) TeacherPage(w http.ResponseWriter, r *http.Request) {
	if !h.authHandler.CheckRole(r, "Преподаватель") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, _ := template.ParseFiles("templates/teacher.html")
	tmpl.Execute(w, map[string]string{"SessionID": r.URL.Query().Get("session_id")})
}

func (h *UIHandler) StudentPage(w http.ResponseWriter, r *http.Request) {
	if !h.authHandler.CheckRole(r, "Студент") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, _ := template.ParseFiles("templates/student.html")
	tmpl.Execute(w, map[string]string{"SessionID": r.URL.Query().Get("session_id")})
}
