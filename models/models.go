package models

type Course struct {
	ID          int
	Title       string
	Description string
	CreditsECTS int
	MaxStudents int
	TeacherID   int
	TeacherName string
}

type Student struct {
	ID               int
	RecordBookNumber string
	EnrollmentYear   int
	Status           string
	LastName         string
	FirstName        string
	MiddleName       string
}

type Grade struct {
	CourseID    int
	CourseTitle string
	Grade       string
	GradeDate   string
}

type Schedule struct {
	ID          int
	DateTime    string
	Room        string
	MeetingLink string
	TeacherID   int
	CourseID    int
	CourseTitle string
	TeacherName string
}
