package controller

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"ie-backend-project/common"
	"ie-backend-project/handler"
	"ie-backend-project/mailer"
	"ie-backend-project/model"
	"net/http"
	"strconv"
	"time"
)

const mailsPerMailer = 5

type Controller struct {
	ch *handler.CourseHandler
	sh *handler.StudentHandler
	sm mailer.Mailer
	v  *validator.Validate
}

func NewController(courseHandler *handler.CourseHandler, studentHandler *handler.StudentHandler, studentMailer mailer.Mailer) *Controller {
	controller := Controller{ch: courseHandler, sh: studentHandler, sm: studentMailer, v: validator.New()}
	return &controller
}

func (h Controller) Register(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.Error{Note: "Function not Implemented"})
}

func (h Controller) Login(c echo.Context) error {
	type Name struct {
		Name string `json:"name" validate:"required"`
		Pass string `json:"pass" validate:"required"`
	}
	data := new(Name)
	bindErr := c.Bind(data)
	parseErr := h.v.Struct(data)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Json. Name shouldn't be empty"})
	}

	//TODO: check password

	claims := &common.JWTCustomClaims{Name: data.Name, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Minute * 5).Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(common.JWTKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.Token{Token: t})
}

func (h Controller) Logout(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.Error{Note: "Function not Implemented"})
}

func (h Controller) NewCourse(c echo.Context) error {
	course := new(model.Course)
	bindErr := c.Bind(course)
	parseErr := h.v.Struct(course)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Course Json, All the course fields are provided"})
	}

	res, err := h.ch.AddCourse(*course)
	if err != nil {
		if errors.Is(err, common.DuplicateCourseError) {
			return c.JSON(http.StatusForbidden, common.Error{Note: "Course Already Exists"})
		}
		return c.JSON(http.StatusExpectationFailed, common.Error{Note: "Couldn't add course"})
	}
	fmt.Println("Added", course)
	return c.JSON(http.StatusCreated, common.ID{ID: res})
}

func (h Controller) GetCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Invalid ID structure. ID should be positive integer"})
	}
	res, err := h.ch.GetCourse(uint(id))
	if errors.Is(err, common.CourseNotFoundError) {
		return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested course"})
	}
	return c.JSON(http.StatusOK, res)
}

func (h Controller) DeleteCourse(c echo.Context) error {
	id := new(common.ID)
	bindErr := c.Bind(id)
	parseErr := h.v.Struct(id)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad ID Json. ID should be a positive integer"})
	}

	if err := h.ch.DeleteCourse(id.ID); errors.Is(err, common.CourseNotFoundError) {
		return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested course"})
	}
	fmt.Println("Deleted Course", id.ID)
	return c.JSON(http.StatusOK, common.Success{Note: "Course deleted"})
}

func (h Controller) GetCourseStudents(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Invalid ID structure. ID should be positive integer"})
	}

	res, err := h.ch.GetStudents(uint(id))
	if err != nil {
		if errors.Is(err, common.CourseNotFoundError) {
			return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested course"})
		}
		if errors.Is(err, common.CourseStudentsError) {
			return c.JSON(http.StatusExpectationFailed, common.Error{Note: "Error while trying to retrieve course students"})
		}
	}
	return c.JSON(http.StatusOK, model.Students{Students: res})
}

func (h Controller) AnnounceCourseResults(c echo.Context) error {
	id := new(common.ID)
	bindErr := c.Bind(id)
	parseErr := h.v.Struct(id)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad ID Json. ID should be a positive integer"})
	}

	res, err := h.ch.GetStudents(id.ID)
	if err != nil {
		if errors.Is(err, common.CourseNotFoundError) {
			return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested course"})
		}
		if errors.Is(err, common.CourseStudentsError) {
			return c.JSON(http.StatusExpectationFailed, common.Error{Note: "Error while trying to retrieve course students"})
		}
	}

	r := common.Results{Results: make([]interface{}, len(res))}
	mailerCount := len(res) / mailsPerMailer
	if mailerCount*mailsPerMailer < len(res) {
		mailerCount++
	}
	//rs := make([]common.Results, mailerCount)
	ch := make(chan int)
	for i := 0; i < mailerCount; i++ {
		beg := i * mailsPerMailer
		end := (i + 1) * mailsPerMailer
		if end > len(res) {
			end = len(res)
		}
		go h.sm.SendMails(res[beg:end], r.Results[beg:end], ch)
		println("Started Mailer", i)
	}
	for i := 0; i < mailerCount; i++ {
		println("Waiting on Mailer", i)
		<-ch
		//r.Results = append(r.Results, (<-ch[i]).Results)
		println("Got the results from Mailer", i)
	}
	return c.JSON(http.StatusOK, r)
}

func (h Controller) UpdateCourseInstructor(c echo.Context) error {
	type CourseInstructor struct {
		ID         uint   `json:"id" validate:"required"`
		Instructor string `json:"instructor" validate:"required"`
	}
	data := new(CourseInstructor)
	bindErr := c.Bind(data)
	parseErr := h.v.Struct(data)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Json. ID should be positive integer"})
	}

	if err := h.ch.UpdateCourseInstructor(data.ID, data.Instructor); err != nil {
		if err == common.InvalidInstructorError {
			return c.JSON(http.StatusBadRequest, common.Error{Note: "Instructor name is required"})
		}
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Course not Found"})
	}
	return c.JSON(http.StatusOK, common.Success{Note: "Course instructor updated"})
}

func (h Controller) NewStudent(c echo.Context) error {
	students := new(model.Students)
	bindErr := c.Bind(students)
	parseErr := h.v.Struct(students)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Students Json. All fields should be filled"})
	}

	r := common.Results{Results: make([]interface{}, 0, len(students.Students))}
	for _, student := range students.Students {
		if parseErr1 := h.v.Struct(student); parseErr1 != nil {
			r.Results = append(r.Results, common.Error{Note: "Not all the student fields are provided"})
			continue
		}

		course, err := h.ch.GetCourse(student.CourseID)
		if errors.Is(err, common.CourseNotFoundError) {
			r.Results = append(r.Results, common.Error{Note: "There is no class with given class_id"})
			continue
		}
		student.Course = *course
		res, err := h.sh.AddStudent(student)
		if err != nil {
			if errors.Is(err, common.DuplicateStudentError) {
				r.Results = append(r.Results, common.Error{Note: "Student already exists"})
				continue
			}
			if errors.Is(err, common.InvalidScoreError) {
				r.Results = append(r.Results, common.Error{Note: "Invalid student score"})
				continue
			}
			if errors.Is(err, common.InvalidEmailError) {
				r.Results = append(r.Results, common.Error{Note: "Invalid email address"})
				continue
			}
			r.Results = append(r.Results, common.Error{Note: "Couldn't add student"})
			continue
		}
		fmt.Println("Added", student)
		r.Results = append(r.Results, common.ID{ID: res})
		continue
	}
	return c.JSON(http.StatusOK, r)
}

func (h Controller) GetStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Invalid ID structure. ID should be positive integer"})
	}
	res, err := h.sh.GetStudent(uint(id))
	if errors.Is(err, common.StudentNotFoundError) || errors.Is(err, common.StudentCourseError) {
		return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested student"})
	}
	return c.JSON(http.StatusOK, res)
}

func (h Controller) UpdateStudentScore(c echo.Context) error {
	type StudentScore struct {
		ID    uint `json:"id" validate:"required"`
		Score uint `json:"score" validate:"gte=0,lte=20"`
	}
	data := new(StudentScore)
	bindErr := c.Bind(data)
	parseErr := h.v.Struct(data)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Json. ID should be positive integer and score should be between 0 an 20"})
	}

	if err := h.sh.UpdateStudentScore(data.ID, data.Score); err != nil {
		if err == common.InvalidScoreError {
			return c.JSON(http.StatusBadRequest, common.Error{Note: "Invalid score value. score should be between 0 an 20"})
		}
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Student not Found"})
	}
	return c.JSON(http.StatusOK, common.Success{Note: "Student score updated"})
}

func (h Controller) UpdateStudentEmail(c echo.Context) error {
	type StudentEmail struct {
		ID    uint   `json:"id" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}
	data := new(StudentEmail)
	bindErr := c.Bind(data)
	parseErr := h.v.Struct(data)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Json. ID should be positive integer and Email should be valid"})
	}

	if err := h.sh.UpdateStudentEmail(data.ID, data.Email); err != nil {
		if err == common.InvalidEmailError {
			return c.JSON(http.StatusBadRequest, common.Error{Note: "Invalid email address"})
		}
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Student not Found"})
	}
	return c.JSON(http.StatusOK, common.Success{Note: "Student email updated"})
}

func (h Controller) DeleteStudent(c echo.Context) error {
	id := new(common.ID)
	bindErr := c.Bind(id)
	parseErr := h.v.Struct(id)
	if bindErr != nil || parseErr != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad ID Json. ID should be a positive integer"})
	}

	if err := h.sh.DeleteStudent(id.ID); errors.Is(err, common.StudentNotFoundError) {
		return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested student"})
	}
	return c.JSON(http.StatusOK, common.Success{Note: "Student Deleted"})
}
