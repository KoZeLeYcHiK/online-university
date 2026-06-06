package services

import (
	"database/sql"
	"online-university/models"
)

type StudentService struct {
	db *sql.DB
}

func NewStudentService(db *sql.DB) *StudentService {
	return &StudentService{db: db}
}

func (s *StudentService) GetAll() ([]models.Student, error) {
	rows, err := s.db.Query(`
        SELECT с.id_студента, с.Номер_зачетки, с.Год_поступления, с.Статус,
               ф.Фамилия, ф.Имя, ф.Отчество
        FROM Студент с
        JOIN Физическое_лицо ф ON с.id_лица = ф.id_лица
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.RecordBookNumber, &s.EnrollmentYear, &s.Status,
			&s.LastName, &s.FirstName, &s.MiddleName)
		if err != nil {
			continue
		}
		students = append(students, s)
	}
	return students, nil
}

func (s *StudentService) Create(lastName, firstName, middleName, recordBook, enrollmentYear, status, directionID string) error {
	var personID int
	err := s.db.QueryRow(`
        INSERT INTO Физическое_лицо (Фамилия, Имя, Отчество, Дата_рождения, Электронная_почта, Телефон)
        VALUES ($1, $2, $3, '2000-01-01', 'temp@temp.ru', '79000000000')
        RETURNING id_лица
    `, lastName, firstName, middleName).Scan(&personID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
        INSERT INTO Студент (Номер_зачетки, Год_поступления, Статус, id_направления, id_лица)
        VALUES ($1, $2, $3, $4, $5)
    `, recordBook, enrollmentYear, status, directionID, personID)
	return err
}

func (s *StudentService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM Студент WHERE id_студента = $1", id)
	return err
}

func (s *StudentService) Update(id int, recordBook, enrollmentYear, status string) error {
	_, err := s.db.Exec(`
        UPDATE Студент 
        SET Номер_зачетки = $1, Год_поступления = $2, Статус = $3
        WHERE id_студента = $4
    `, recordBook, enrollmentYear, status, id)
	return err
}
