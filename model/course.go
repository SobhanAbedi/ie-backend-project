package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name       string `json:"name" validate:"required"`
	Instructor string `json:"instructor" validate:"required"`
}

func (c Course) String() string {
	return fmt.Sprintf("%s Course by %s", c.Name, c.Instructor)
}
