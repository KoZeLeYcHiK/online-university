package main

import (
	"fmt"
	"log"
	"net/http"
	"online-university/database"
	"online-university/handlers"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	fmt.Println("Подключение к PostgreSQL установлено")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Страницы
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	http.HandleFunc("/admin", handlers.AdminPage)
	http.HandleFunc("/teacher", handlers.TeacherPage)
	http.HandleFunc("/student", handlers.StudentPage)

	// API - общие
	http.HandleFunc("/api/courses", handlers.GetCoursesHandler)
	http.HandleFunc("/api/students", handlers.GetStudentsHandler)
	http.HandleFunc("/api/teachers", handlers.GetTeachersHandler)
	http.HandleFunc("/api/schedule", handlers.GetScheduleHandler)
	http.HandleFunc("/api/students-grades", handlers.GetAllGradesHandler)

	// Admin управление
	http.HandleFunc("/api/add-student", handlers.AddStudentHandler)
	http.HandleFunc("/api/delete-student", handlers.DeleteStudentHandler)
	http.HandleFunc("/api/add-course", handlers.AddCourseHandler)
	http.HandleFunc("/api/delete-course", handlers.DeleteCourseHandler)
	http.HandleFunc("/api/add-teacher", handlers.AddTeacherHandler)
	http.HandleFunc("/api/delete-teacher", handlers.DeleteTeacherHandler)

	// API - преподаватель
	http.HandleFunc("/api/teacher/courses", handlers.GetTeacherCoursesHandler)
	http.HandleFunc("/api/teacher/update-grade", handlers.UpdateGradeHandler)
	http.HandleFunc("/api/course/students", handlers.GetCourseStudentsHandler)

	// API - студент
	http.HandleFunc("/api/student/grades", handlers.GetStudentGradesHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
