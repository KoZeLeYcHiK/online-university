function loadPage(page) {
    const content = document.getElementById('content');
    content.innerHTML = '<p>Загрузка...</p>';
    
    let apiUrl = '';
    let title = '';
    switch(page) {
        case 'courses':
            apiUrl = '/api/courses';
            title = 'Курсы';
            break;
        case 'students':
            apiUrl = '/api/students';
            title = 'Студенты';
            break;
        case 'teachers':
            apiUrl = '/api/teachers';
            title = 'Преподаватели';
            break;
        case 'schedule':
            apiUrl = '/api/schedule';
            title = 'Расписание';
            break;
        case 'studentsGrades':
            apiUrl = '/api/students-grades';
            title = 'Успеваемость студентов';
            break;
        case 'myCourses':
            apiUrl = '/api/teacher/courses';
            title = 'Мои курсы';
            break;
        case 'myGrades':
            apiUrl = '/api/student/grades';
            title = 'Моя успеваемость';
            break;
        default:
            apiUrl = '/api/' + page;
            title = page;
    }
    
    fetch(apiUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error('HTTP error ' + response.status);
            }
            return response.json();
        })
        .then(data => {
            if (!data || data.length === 0) {
                content.innerHTML = `<h2>${title}</h2><p>Нет данных</p>`;
                return;
            }
            
            let headers = {};
            if (page === 'studentsGrades') {
                headers = { 'student': 'Студент', 'course': 'Курс', 'grade': 'Оценка', 'date': 'Дата' };
            } else if (page === 'students') {
                headers = { 'ID': 'ID', 'RecordBookNumber': 'Номер зачетки', 'EnrollmentYear': 'Год поступления', 'Status': 'Статус', 'LastName': 'Фамилия', 'FirstName': 'Имя', 'MiddleName': 'Отчество' };
            } else if (page === 'courses') {
                headers = { 'ID': 'ID', 'Title': 'Название', 'Description': 'Описание', 'CreditsECTS': 'Кредиты', 'MaxStudents': 'Макс. студентов', 'TeacherName': 'Преподаватель' };
            } else if (page === 'teachers') {
                headers = { 'full_name': 'ФИО', 'department': 'Кафедра', 'qualification': 'Квалификация', 'degree': 'Ученая степень' };
            } else if (page === 'schedule') {
                headers = { 'course': 'Курс', 'teacher': 'Преподаватель', 'datetime': 'Дата и время', 'room': 'Аудитория', 'meeting_link': 'Ссылка' };
            } else if (page === 'myGrades') {
                headers = { 'course': 'Курс', 'grade': 'Оценка', 'date': 'Дата' };
            } else if (page === 'myCourses') {
                headers = { 'id': 'ID', 'title': 'Название', 'description': 'Описание', 'credits_ects': 'Кредиты', 'max_students': 'Макс. студентов' };
                let html = `<h2>${title}</h2><div class="table-container"><table><thead><tr>`;
                for (let key in headers) {
                    html += `<th>${headers[key]}</th>`;
                }
                html += '<th>Действие</th></thead><tbody>';
                data.forEach(item => {
                    html += '<tr>';
                    for (let key in headers) {
                        let value = item[key];
                        if (value === null || value === undefined) value = '-';
                        html += `<td>${value}</td>`;
                    }
                    html += `<td><button class="btn btn-primary" onclick="showCourseStudents(${item.id}, '${item.title}')">Список студентов</button></td>`;
                    html += '</tr>';
                });
                html += '</tbody></table></div>';
                content.innerHTML = html;
                return;
            } else {
                Object.keys(data[0]).forEach(key => {
                    let headerName = key;
                    if (key === 'id') headerName = 'ID';
                    else if (key === 'title') headerName = 'Название';
                    else if (key === 'full_name') headerName = 'ФИО';
                    else if (key === 'grade') headerName = 'Оценка';
                    else if (key === 'date') headerName = 'Дата';
                    else if (key === 'student') headerName = 'Студент';
                    else if (key === 'course') headerName = 'Курс';
                    else headerName = key;
                    headers[key] = headerName;
                });
            }
            
            let html = `<h2>${title}</h2><div class="table-container"><table><thead><tr>`;
            for (let key in headers) {
                html += `<th>${headers[key]}</th>`;
            }
            html += '</td></thead><tbody>';
            data.forEach(item => {
                html += '<tr>';
                for (let key in headers) {
                    let value = item[key];
                    if (value === null || value === undefined) value = '-';
                    html += `<td>${value}</td>`;
                }
                html += '</tr>';
            });
            html += '</tbody></table></div>';
            content.innerHTML = html;
        })
        .catch((error) => {
            console.error('Ошибка:', error);
            content.innerHTML = `<h2>${title}</h2><p style="color:red">Ошибка загрузки данных</p>`;
        });
}

function showCourseStudents(courseId, courseTitle) {
    const content = document.getElementById('content');
    content.innerHTML = '<p>Загрузка...</p>';
    
    fetch(`/api/course/students?course_id=${courseId}`)
        .then(response => response.json())
        .then(data => {
            if (!data || data.length === 0) {
                content.innerHTML = `<h2>${courseTitle}</h2><p>Нет студентов на курсе</p><button onclick="loadPage('myCourses')">Назад</button>`;
                return;
            }
            
            let html = `<h2>${courseTitle} - Студенты</h2><div class="table-container"><table><thead><tr>`;
            html += '<th>Студент</th><th>Оценка</th><th>Действие</th></thead><tbody>';
            data.forEach(student => {
                html += `<tr>
                    <td>${student.full_name}</td>
                    <td>${student.grade || 'не оценено'}</td>
                    <td><button class="btn btn-primary" onclick="openGradeModal(${student.id}, ${courseId}, '${student.full_name}')">Выставить оценку</button></td>
                </tr>`;
            });
            html += '</tbody></table><br><button class="btn" onclick="loadPage(\'myCourses\')">Назад к курсам</button></div>';
            content.innerHTML = html;
        })
        .catch(() => {
            content.innerHTML = `<h2>${courseTitle}</h2><p style="color:red">Ошибка загрузки студентов</p><button onclick="loadPage('myCourses')">Назад</button>`;
        });
}

function openGradeModal(studentId, courseId, studentName) {
    document.getElementById('gradeStudentId').value = studentId;
    document.getElementById('gradeCourseId').value = courseId;
    document.getElementById('studentName').innerText = studentName;
    document.getElementById('gradeModal').style.display = 'block';
}

function closeGradeModal() {
    document.getElementById('gradeModal').style.display = 'none';
}

document.addEventListener('DOMContentLoaded', function() {
    const gradeForm = document.getElementById('gradeForm');
    if (gradeForm) {
        gradeForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const studentId = document.getElementById('gradeStudentId').value;
            const courseId = document.getElementById('gradeCourseId').value;
            const grade = document.querySelector('#gradeForm select[name="grade"]').value;
            
            const formData = new FormData();
            formData.append('student_id', studentId);
            formData.append('course_id', courseId);
            formData.append('grade', grade);
            
            fetch('/api/teacher/update-grade', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    alert('Оценка выставлена успешно!');
                    closeGradeModal();
                    const courseTitle = document.querySelector('#gradeModal h3').innerText;
                    showCourseStudents(courseId, courseTitle);
                } else {
                    alert('Ошибка при выставлении оценки');
                }
            })
            .catch(() => {
                alert('Ошибка при выставлении оценки');
            });
        });
    }
});