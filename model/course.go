package model

import (
	"encoding/xml"
	"fmt"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	XMLName    xml.Name `xml:"course" validate:"-" gorm:"-:all"`
	Name       string   `json:"name" xml:"name" validate:"required"`
	Instructor string   `json:"instructor" xml:"instructor" validate:"required"`
}

func (c Course) String() string {
	return fmt.Sprintf("%s Course by %s", c.Name, c.Instructor)
}

type Courses struct {
	XMLName xml.Name `xml:"courses" validate:"-" gorm:"-:all"`
	Courses []Course `json:"courses" xml:"course" validate:"required"`
}
