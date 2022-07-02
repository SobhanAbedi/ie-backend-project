package handler

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ie-backend-project/common"
	"ie-backend-project/model"
	"net/mail"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(dsn string) (*StudentHandler, error) {
	db, err := gorm.Open(sqlite.Open("db/"+dsn), &gorm.Config{})
	if err != nil {
		return nil, common.DBConnectionFailedError
	}

	err = db.AutoMigrate(&model.Student{})
	if err != nil {
		return nil, common.StudentMMFailedError
	}

	handler := StudentHandler{db}
	return &handler, nil
}

func (h StudentHandler) NewStudent(firstName, lastName, email string, score int, course *model.Course) (*model.Student, error) {
	if score < 0 || score > 20 {
		return nil, common.InvalidScoreError
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, common.InvalidEmailError
	}
	newStd := model.Student{FirstName: firstName, LastName: lastName, Email: email, Score: score, CourseID: course.ID, Course: *course}
	if h.Exists(newStd) {
		return nil, common.DuplicateStudentError
	}
	res := h.db.Create(&newStd)
	if res.Error != nil {
		return nil, res.Error
	}
	return &newStd, nil
}

func (h StudentHandler) AddStudent(student model.Student) (uint, error) {
	if student.Score < 0 || student.Score > 20 {
		return 0, common.InvalidScoreError
	}
	if _, err := mail.ParseAddress(student.Email); err != nil {
		return 0, common.InvalidEmailError
	}
	if h.Exists(student) {
		return 0, common.DuplicateStudentError
	}
	res := h.db.Create(&student)
	if res.Error != nil {
		return 0, res.Error
	}
	return student.ID, nil
}

func (h StudentHandler) GetStudent(id uint) (*model.Student, error) {
	std := new(model.Student)
	h.db.Limit(1).Find(std, id)
	if std.ID == 0 {
		return nil, common.StudentNotFoundError
	}
	course := new(model.Course)
	h.db.Limit(1).Find(course, std.CourseID)
	if course.ID == 0 {
		h.db.Delete(&model.Student{}, std.ID)
		fmt.Println("Deleted Student", std.ID)
		return nil, common.StudentClassError
	}
	std.Course = *course
	return std, nil
}

func (h StudentHandler) UpdateStudentScore(id, newScore uint) error {
	if newScore < 0 || newScore > 20 {
		return common.InvalidScoreError
	}
	std, err := h.GetStudent(id)
	if err != nil {
		return err
	}
	h.db.Model(std).Update("Score", newScore)
	return nil
}

func (h StudentHandler) UpdateStudentEmail(id uint, newEmail string) error {
	if _, err := mail.ParseAddress(newEmail); err != nil {
		return common.InvalidEmailError
	}
	std, err := h.GetStudent(id)
	if err != nil {
		return err
	}
	h.db.Model(std).Update("Email", newEmail)
	return nil
}

func (h StudentHandler) DeleteStudent(id uint) error {
	student := new(model.Student)
	h.db.Limit(1).Find(student, id)
	if student.ID == 0 {
		return common.StudentNotFoundError
	}
	h.db.Delete(&model.Student{}, id)
	return nil
}

func (h StudentHandler) Exists(student model.Student) bool {
	foundOne := new(model.Student)
	h.db.Where(&model.Student{FirstName: student.FirstName, LastName: student.LastName, CourseID: student.CourseID}).Limit(1).Find(foundOne)
	if foundOne.ID == 0 {
		return false
	}
	studentCourse := new(model.Course)
	h.db.Limit(1).Find(studentCourse, foundOne.CourseID)
	if studentCourse.ID == 0 {
		h.db.Delete(&model.Student{}, foundOne.ID)
		fmt.Println("Deleted Student", foundOne.ID)
		return false
	}
	return true
}
