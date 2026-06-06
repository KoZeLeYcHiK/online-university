package services

import (
	"database/sql"
)

type TeacherService struct {
	db *sql.DB
}

func NewTeacherService(db *sql.DB) *TeacherService {
	return &TeacherService{db: db}
}

func (s *TeacherService) GetAll() ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
        SELECT ф.Фамилия, ф.Имя, ф.Отчество, п.Кафедра, кв.Название, уч.Название
        FROM Преподаватель п
        JOIN Физическое_лицо ф ON п.id_лица = ф.id_лица
        JOIN Квалификация кв ON п.id_квалификации = кв.id_квалификации
        JOIN Ученая_степень уч ON п.id_ученой_степени = уч.id_степени
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []map[string]interface{}
	for rows.Next() {
		var lastName, firstName, middleName, department, qualification, degree string
		err := rows.Scan(&lastName, &firstName, &middleName, &department, &qualification, &degree)
		if err != nil {
			continue
		}
		teachers = append(teachers, map[string]interface{}{
			"full_name":     lastName + " " + firstName + " " + middleName,
			"department":    department,
			"qualification": qualification,
			"degree":        degree,
		})
	}
	return teachers, nil
}

func (s *TeacherService) Create(lastName, firstName, middleName, department, qualificationID, degreeID string) error {
	var personID int
	err := s.db.QueryRow(`
        INSERT INTO Физическое_лицо (Фамилия, Имя, Отчество, Дата_рождения, Электронная_почта, Телефон)
        VALUES ($1, $2, $3, '1980-01-01', 'temp@temp.ru', '79000000000')
        RETURNING id_лица
    `, lastName, firstName, middleName).Scan(&personID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
        INSERT INTO Преподаватель (Кафедра, id_квалификации, id_ученой_степени, id_лица)
        VALUES ($1, $2, $3, $4)
    `, department, qualificationID, degreeID, personID)
	return err
}

func (s *TeacherService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM Преподаватель WHERE id_преподавателя = $1", id)
	return err
}

func (s *TeacherService) Update(id int, department, qualificationID, degreeID string) error {
	_, err := s.db.Exec(`
        UPDATE Преподаватель 
        SET Кафедра = $1, id_квалификации = $2, id_ученой_степени = $3
        WHERE id_преподавателя = $4
    `, department, qualificationID, degreeID, id)
	return err
}
