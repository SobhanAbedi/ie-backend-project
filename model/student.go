package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Score     int    `json:"score"`
	CourseID  uint   `json:"course_id"`
	Course    Course `json:"course"`
}

func (s Student) String() string {
	return fmt.Sprintf("%s %s scored %d in %s course by %s", s.FirstName, s.LastName, s.Score, s.Course.Name, s.Course.Instructor)
}

type Students struct {
	Students []Student `json:"students"`
}
