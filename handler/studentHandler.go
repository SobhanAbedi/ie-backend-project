package handler

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ie-backend-project/model"
	"net/mail"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(dsn string) (*StudentHandler, error) {
	db, err := gorm.Open(sqlite.Open("db/"+dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("database connection failed")
	}

	err = db.AutoMigrate(&model.Student{})
	if err != nil {
		return nil, errors.New("student model migration failed")
	}

	handler := StudentHandler{db}
	return &handler, nil
}

func (h StudentHandler) NewStudent(firstName, lastName, email string, score int, course *model.Course) (*model.Student, error) {
	newStd := model.Student{FirstName: firstName, LastName: lastName, Email: email, Score: score, Course: *course}
	res := h.db.Create(&newStd)
	if res.Error != nil {
		return nil, res.Error
	}
	return &newStd, nil
}

func (h StudentHandler) GetStudent(id uint) (*model.Student, error) {
	std := model.Student{}
	res := h.db.First(&std, id)
	if res.Error != nil {
		return nil, res.Error
	}
	course := model.Course{}
	res = h.db.First(&course, std.CourseID)
	if res.Error != nil {
		return nil, res.Error
	}
	std.Course = course
	return &std, nil
}

func (h StudentHandler) UpdateStudentScore(id uint, newScore int) error {
	if newScore < 0 || newScore > 20 {
		return errors.New("invalid score")
	}
	h.db.Model(&model.Student{}).Where("ID = ?", id).Update("Score", newScore)
	return nil
}

func (h StudentHandler) UpdateStudentEmail(id uint, newEmail string) error {
	_, err := mail.ParseAddress(newEmail)
	if err != nil {
		return err
	}
	h.db.Model(&model.Student{}).Where("ID = ?", id).Update("Email", newEmail)
	return nil
}

func (h StudentHandler) DeleteStudent(id uint) {
	h.db.Delete(&model.Student{}, id)
}
