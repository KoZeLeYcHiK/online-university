package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"online-university/database"
)

// Хранилище сессий (в памяти)
var sessions = make(map[string]string)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	errorMsg := r.URL.Query().Get("error")
	tmpl.Execute(w, map[string]interface{}{"error": errorMsg == "1"})
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== AuthHandler вызван ===")

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	fmt.Println("Login:", login)
	fmt.Println("Password:", password)

	var role string
	var userID int
	err := database.DB.QueryRow(`
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

	// Создаём ID сессии
	sessionID := fmt.Sprintf("%d_%s", userID, login)
	sessions[sessionID] = role

	// Редирект с session_id
	switch role {
	case "Администратор":
		http.Redirect(w, r, "/admin?session_id="+sessionID, http.StatusSeeOther)
	case "Преподаватель":
		http.Redirect(w, r, "/teacher?session_id="+sessionID, http.StatusSeeOther)
	case "Студент":
		http.Redirect(w, r, "/student?session_id="+sessionID, http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	delete(sessions, sessionID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func GetRole(r *http.Request) string {
	sessionID := r.URL.Query().Get("session_id")
	role, exists := sessions[sessionID]
	if !exists {
		return ""
	}
	return role
}

func CheckRole(r *http.Request, allowedRoles ...string) bool {
	role := GetRole(r)
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}
