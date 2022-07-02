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
	CourseID  int    `json:"course"`
	Course    Course
}

func (s Student) String() string {
	return fmt.Sprintf("%s %s from %s class witdh a score of %d", s.FirstName, s.LastName, s.Course.Name, s.Score)
}
