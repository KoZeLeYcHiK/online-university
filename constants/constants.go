package constants

// Роли пользователей
const (
	RoleAdmin   = "Администратор"
	RoleTeacher = "Преподаватель"
	RoleStudent = "Студент"
)

// Оценки
const (
	GradeExcellent = "отлично"
	GradeGood      = "хорошо"
	GradeSatisf    = "удовлетворительно"
	GradeUnsatisf  = "неудовлетворительно"
)

// Статусы студентов
const (
	StatusActive        = "активен"
	StatusAcademicLeave = "академический отпуск"
	StatusExpelled      = "отчислен"
)

// Уровни доступа
const (
	AccessLevelFull  = "полный"
	AccessLevelLimit = "ограниченный"
)

// Статусы курсов (если есть)
const (
	CourseStatusActive = "активен"
	CourseStatusClosed = "закрыт"
)
