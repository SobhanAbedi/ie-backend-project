package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Score     int    `json:"score" validate:"gte=0,lte=20"`
	CourseID  uint   `json:"course_id" validate:"required"`
	Course    Course `json:"course" validate:"-"`
}

func (s Student) String() string {
	return fmt.Sprintf("%s %s scored %d in %s course by %s", s.FirstName, s.LastName, s.Score, s.Course.Name, s.Course.Instructor)
}

type Students struct {
	Students []Student `json:"students" validate:"required"`
}
