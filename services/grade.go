package services

import (
	"database/sql"
)

type GradeService struct {
	db *sql.DB
}

func NewGradeService(db *sql.DB) *GradeService {
	return &GradeService{db: db}
}

func (s *GradeService) GetStudentGrades(studentID int) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
        SELECT к.Название, COALESCE(у.Оценка, 'не оценено'), to_char(у.Дата_оценки, 'DD.MM.YYYY')
        FROM Успеваемость у
        JOIN Курс к ON у.id_курса = к.id_курса
        WHERE у.id_студента = $1
    `, studentID)
	if err != nil {
		return nil, err
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
	return grades, nil
}

func (s *GradeService) GetAllGrades() ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
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
		return nil, err
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
	return grades, nil
}

func (s *GradeService) UpdateGrade(studentID, courseID, grade string) error {
	_, err := s.db.Exec(`
        UPDATE Успеваемость 
        SET Оценка = $1, Дата_оценки = CURRENT_DATE
        WHERE id_студента = $2 AND id_курса = $3
    `, grade, studentID, courseID)
	return err
}

func (s *GradeService) GetCourseStudents(courseID string) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
        SELECT с.id_студента, ф.Фамилия, ф.Имя, ф.Отчество, COALESCE(у.Оценка, '')
        FROM Успеваемость у
        JOIN Студент с ON у.id_студента = с.id_студента
        JOIN Физическое_лицо ф ON с.id_лица = ф.id_лица
        WHERE у.id_курса = $1
    `, courseID)
	if err != nil {
		return nil, err
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
	return students, nil
}
