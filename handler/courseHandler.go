package handler

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ie-backend-project/common"
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
	if h.Exists(newCourse) {
		return nil, errors.New("duplicate course")
	}
	res := h.db.Create(&newCourse)
	if res.Error != nil {
		return nil, res.Error
	}
	return &newCourse, nil
}

func (h CourseHandler) AddCourse(course model.Course) (uint, error) {
	if h.Exists(course) {
		return 0, common.DuplicateCourseError
	}
	res := h.db.Create(&course)
	if res.Error != nil {
		return 0, res.Error
	}
	return course.ID, nil
}

func (h CourseHandler) GetCourse(id uint) (*model.Course, error) {
	course := new(model.Course)
	h.db.Limit(1).Find(course, id)
	if course.ID != 0 {
		return course, nil
	}
	return nil, common.CourseNotFoundError

}

func (h CourseHandler) DeleteCourse(id uint) error {
	course := new(model.Course)
	h.db.Limit(1).Find(course, id)
	if course.ID == 0 {
		return common.CourseNotFoundError
	}
	h.db.Delete(&model.Course{}, id)
	return nil
}

func (h CourseHandler) Exists(course model.Course) bool {
	foundOne := new(model.Course)
	h.db.Where(&model.Course{Name: course.Name, Instructor: course.Instructor}).Limit(1).Find(foundOne)
	if foundOne.ID != 0 {
		return true
	}
	return false
}
