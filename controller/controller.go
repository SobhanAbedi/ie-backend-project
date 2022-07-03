package controller

import (
	"errors"
	"fmt"
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
}

func NewController(courseHandler *handler.CourseHandler, studentHandler *handler.StudentHandler, studentMailer mailer.Mailer) *Controller {
	controller := Controller{ch: courseHandler, sh: studentHandler, sm: studentMailer}
	return &controller
}

func (h Controller) Register(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.Error{Note: "Function not Implemented"})
}

func (h Controller) Login(c echo.Context) error {
	type Name struct {
		Name string `json:"name"`
		Pass string `json:"password"`
	}
	data := new(Name)
	if err := c.Bind(data); err != nil || data.Name == "" {
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
	if err := c.Bind(course); err != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Course Json"})
	}

	if course.Name == "" || course.Instructor == "" {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Not all the course fields are provided"})
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
	if err := c.Bind(id); err != nil || id.ID == 0 {
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
	if err := c.Bind(id); err != nil || id.ID == 0 {
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
		ID         uint   `json:"id"`
		Instructor string `json:"instructor"`
	}
	data := new(CourseInstructor)
	if err := c.Bind(data); err != nil || data.ID == 0 || data.Instructor == "" {
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
	if err := c.Bind(students); err != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Students Json"})
	}

	r := common.Results{Results: make([]interface{}, 0, len(students.Students))}
	for _, student := range students.Students {
		if student.FirstName == "" || student.LastName == "" || student.Email == "" || student.CourseID == 0 {
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
		ID    uint `json:"id"`
		Score uint `json:"score"`
	}
	data := new(StudentScore)
	if err := c.Bind(data); err != nil || data.ID == 0 {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Json. ID should be positive integer"})
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
		ID    uint   `json:"id"`
		Email string `json:"email"`
	}
	data := new(StudentEmail)
	if err := c.Bind(data); err != nil || data.ID == 0 {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Bad Json. ID should be positive integer"})
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
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.Error{Note: "Invalid ID structure. ID should be positive integer"})
	}
	if err := h.sh.DeleteStudent(uint(id)); errors.Is(err, common.StudentNotFoundError) {
		return c.JSON(http.StatusNotFound, common.Error{Note: "Couldn't find requested student"})
	}
	return c.JSON(http.StatusOK, common.Success{Note: "Student Deleted"})
}
