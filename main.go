package main

import (
	"fmt"
	"log"
	"net/http"
	"online-university/config"
	"online-university/database"
	"online-university/handlers"
	"online-university/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()
	fmt.Println("Подключение к PostgreSQL установлено")

	// Сервисы
	courseService := services.NewCourseService(db.Conn)
	studentService := services.NewStudentService(db.Conn)
	teacherService := services.NewTeacherService(db.Conn)
	gradeService := services.NewGradeService(db.Conn)

	// Обработчики
	authHandler := handlers.NewAuthHandler(db.Conn)
	courseHandler := handlers.NewCourseHandler(courseService, authHandler)
	studentHandler := handlers.NewStudentHandler(studentService, authHandler)
	teacherHandler := handlers.NewTeacherHandler(teacherService, authHandler)
	gradeHandler := handlers.NewGradeHandler(gradeService, authHandler)
	scheduleHandler := handlers.NewScheduleHandler(db.Conn, authHandler)
	uiHandler := handlers.NewUIHandler(db.Conn, authHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Страницы
	http.HandleFunc("/login", authHandler.LoginPage)
	http.HandleFunc("/auth", authHandler.AuthHandler)
	http.HandleFunc("/logout", authHandler.LogoutHandler)

	http.HandleFunc("/admin", uiHandler.AdminPage)
	http.HandleFunc("/teacher", uiHandler.TeacherPage)
	http.HandleFunc("/student", uiHandler.StudentPage)

	// API
	http.HandleFunc("/api/courses", courseHandler.GetCourses)
	http.HandleFunc("/api/students", studentHandler.GetStudents)
	http.HandleFunc("/api/teachers", teacherHandler.GetTeachers)
	http.HandleFunc("/api/schedule", scheduleHandler.GetSchedule)
	http.HandleFunc("/api/students-grades", gradeHandler.GetAllGrades)
	http.HandleFunc("/api/teacher/courses", courseHandler.GetTeacherCourses)
	http.HandleFunc("/api/teacher/update-grade", gradeHandler.UpdateGrade)
	http.HandleFunc("/api/course/students", gradeHandler.GetCourseStudents)
	http.HandleFunc("/api/student/grades", gradeHandler.GetStudentGrades)

	// CRUD администратора
	http.HandleFunc("/api/add-student", studentHandler.AddStudent)
	http.HandleFunc("/api/delete-student", studentHandler.DeleteStudent)
	http.HandleFunc("/api/update-student", studentHandler.UpdateStudent)

	http.HandleFunc("/api/add-course", courseHandler.AddCourse)
	http.HandleFunc("/api/delete-course", courseHandler.DeleteCourse)
	http.HandleFunc("/api/update-course", courseHandler.UpdateCourse)

	http.HandleFunc("/api/add-teacher", teacherHandler.AddTeacher)
	http.HandleFunc("/api/delete-teacher", teacherHandler.DeleteTeacher)
	http.HandleFunc("/api/update-teacher", teacherHandler.UpdateTeacher)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
