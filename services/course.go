package services

import (
	"database/sql"
	"online-university/models"
)

type CourseService struct {
	db *sql.DB
}

func NewCourseService(db *sql.DB) *CourseService {
	return &CourseService{db: db}
}

func (s *CourseService) GetAll() ([]models.Course, error) {
	rows, err := s.db.Query(`
        SELECT к.id_курса, к.Название, к.Описание, к.Кредиты_ECTS, к.Макс_студентов,
               ф.Фамилия, ф.Имя, ф.Отчество
        FROM Курс к
        JOIN Преподаватель п ON к.id_преподавателя = п.id_преподавателя
        JOIN Физическое_лицо ф ON п.id_лица = ф.id_лица
    `)
	if err != nil {
		return nil, err
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
	return courses, nil
}

func (s *CourseService) GetTeacherCourses(teacherID int) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
        SELECT id_курса, Название, Описание, Кредиты_ECTS, Макс_студентов
        FROM Курс
        WHERE id_преподавателя = $1
    `, teacherID)
	if err != nil {
		return nil, err
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
	return courses, nil
}

func (s *CourseService) Create(title, description string, credits, maxStudents, teacherID int) error {
	_, err := s.db.Exec(`
        INSERT INTO Курс (Название, Описание, Кредиты_ECTS, Макс_студентов, id_преподавателя)
        VALUES ($1, $2, $3, $4, $5)
    `, title, description, credits, maxStudents, teacherID)
	return err
}

func (s *CourseService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM Курс WHERE id_курса = $1", id)
	return err
}

func (s *CourseService) Update(id int, title, description string, credits, maxStudents, teacherID int) error {
	_, err := s.db.Exec(`
        UPDATE Курс 
        SET Название = $1, Описание = $2, Кредиты_ECTS = $3, Макс_студентов = $4, id_преподавателя = $5
        WHERE id_курса = $6
    `, title, description, credits, maxStudents, teacherID, id)
	return err
}
