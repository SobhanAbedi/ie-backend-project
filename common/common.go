package common

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

var DBConnectionFailedError = errors.New("database connection failed")

var CourseMMFailedError = errors.New("course model migration failed")
var DuplicateCourseError = errors.New("duplicate course")
var CourseNotFoundError = errors.New("course not found")
var InvalidInstructorError = errors.New("instructor field is empty")
var CourseStudentsError = errors.New("error while trying to retrieve course students")

var StudentMMFailedError = errors.New("student model migration failed")
var DuplicateStudentError = errors.New("duplicate student")
var StudentNotFoundError = errors.New("student not found")
var InvalidScoreError = errors.New("invalid score")
var InvalidEmailError = errors.New("invalid email")
var StudentCourseError = errors.New("invalid class id for student")

type Error struct {
	Note string `json:"error"`
}

type Success struct {
	Note string `json:"msg"`
}

type Token struct {
	Token string `json:"token"`
}

type ID struct {
	ID uint `json:"id" validate:"required"`
}

type Results struct {
	Results []interface{} `json:"results"`
}

type JWTCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

var JWTKey = []byte("my_surprisingly_secret_key")
