package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"online-university/constants"
)

type AuthHandler struct {
	db       *sql.DB
	sessions map[string]string
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		db:       db,
		sessions: make(map[string]string),
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	errorMsg := r.URL.Query().Get("error")
	tmpl.Execute(w, map[string]interface{}{"error": errorMsg == "1"})
}

func (h *AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== AuthHandler вызван ===")

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	var role string
	var userID int
	err := h.db.QueryRow(`
        SELECT id_пользователя, р.Название 
        FROM Пользователь п
        JOIN Роль р ON п.id_роли = р.id_роли
        WHERE Логин = $1 AND Пароль = $2
    `, login, password).Scan(&userID, &role)

	if err != nil {
		fmt.Println("Ошибка аутентификации:", err)
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	fmt.Println("Успешный вход! Роль:", role)

	sessionID := fmt.Sprintf("%d_%s", userID, login)
	h.sessions[sessionID] = role

	switch role {
	case constants.RoleAdmin:
		http.Redirect(w, r, "/admin?session_id="+sessionID, http.StatusSeeOther)
	case constants.RoleTeacher:
		http.Redirect(w, r, "/teacher?session_id="+sessionID, http.StatusSeeOther)
	case constants.RoleStudent:
		http.Redirect(w, r, "/student?session_id="+sessionID, http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	delete(h.sessions, sessionID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) GetRole(r *http.Request) string {
	sessionID := r.URL.Query().Get("session_id")
	role, exists := h.sessions[sessionID]
	if !exists {
		return ""
	}
	return role
}

func (h *AuthHandler) CheckRole(r *http.Request, allowedRoles ...string) bool {
	role := h.GetRole(r)
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}
