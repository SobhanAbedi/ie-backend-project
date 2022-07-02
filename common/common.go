package common

import "errors"

var DuplicateCourseError = errors.New("duplicate course")
var CourseNotFoundError = errors.New("course not found")

type Error struct {
	Note string `json:"error"`
}

type Created struct {
	ID uint `json:"id"`
}
