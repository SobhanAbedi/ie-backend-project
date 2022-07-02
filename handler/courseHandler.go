package handler

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ie-backend-project/model"
)

type CourseHandler struct {
	db *gorm.DB
}

func NewCourseHandler(dsn string) (*CourseHandler, error) {
	db, err := gorm.Open(sqlite.Open("db/"+dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("database connection failed")
	}

	err = db.AutoMigrate(&model.Course{})
	if err != nil {
		return nil, errors.New("course model migration failed")
	}

	handler := CourseHandler{db}
	return &handler, nil
}

func (h CourseHandler) NewCourse(name, instructor string) (*model.Course, error) {
	newCourse := model.Course{Name: name, Instructor: instructor}
	res := h.db.Create(&newCourse)
	if res.Error != nil {
		return nil, res.Error
	}
	return &newCourse, nil

}

func (h CourseHandler) GetCourse(id uint) (*model.Course, error) {
	course := model.Course{}
	res := h.db.First(&course, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &course, nil
}

func (h CourseHandler) DeleteCourse(id uint) {
	h.db.Delete(&model.Course{}, id)
}
