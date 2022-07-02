package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"ie-backend-project/common"
	"ie-backend-project/handler"
	"ie-backend-project/model"
	"net/http"
	"strconv"
)

func remErr[T any](output T, err error) T {
	if err != nil {
		fmt.Println(err)
	}
	return output
}

type Controller struct {
	ch *handler.CourseHandler
	sh *handler.StudentHandler
}

func NewController(courseHandler *handler.CourseHandler, studentHandler *handler.StudentHandler) *Controller {
	controller := Controller{ch: courseHandler, sh: studentHandler}
	return &controller
}

func (h Controller) Register(c echo.Context) error {
	return nil
}

func (h Controller) Login(c echo.Context) error {
	return nil
}

func (h Controller) Logout(c echo.Context) error {
	return nil
}

func (h Controller) NewCourse(c echo.Context) error {
	course := new(model.Course)
	if err := c.Bind(course); err != nil {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Bad Course Json"}))))
	}

	if course.Name == "" || course.Instructor == "" {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Not all the course fields are provided"}))))
	}

	res, err := h.ch.AddCourse(*course)
	if err != nil {
		if errors.Is(err, common.DuplicateCourseError) {
			return c.String(http.StatusForbidden, string(remErr(json.Marshal(common.Error{Note: "Course Already Exists"}))))
		}
		return c.String(http.StatusExpectationFailed, string(remErr(json.Marshal(common.Error{Note: "Couldn't add course"}))))
	}
	fmt.Println("Added", course)
	return c.String(http.StatusCreated, string(remErr(json.Marshal(common.ID{ID: res}))))
}

func (h Controller) GetCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Invalid ID structure. ID should be positive integer"}))))
	}
	res, err := h.ch.GetCourse(uint(id))
	if errors.Is(err, common.CourseNotFoundError) {
		return c.String(http.StatusNotFound, string(remErr(json.Marshal(common.Error{Note: "Couldn't find requested course"}))))
	}
	return c.String(http.StatusOK, string(remErr(json.Marshal(res))))
}

func (h Controller) DeleteCourse(c echo.Context) error {
	id := new(common.ID)
	if err := c.Bind(id); err != nil || id.ID == 0 {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Bad ID Json. ID should be a positive integer"}))))
	}

	if err := h.ch.DeleteCourse(id.ID); errors.Is(err, common.CourseNotFoundError) {
		return c.String(http.StatusNotFound, string(remErr(json.Marshal(common.Error{Note: "Couldn't find requested course"}))))
	}
	fmt.Println("Deleted Course", id.ID)
	return c.String(http.StatusOK, string(remErr(json.Marshal(common.Success{Note: "Course deleted"}))))
}

func (h Controller) NewStudent(c echo.Context) error {
	students := new(model.Students)
	if err := c.Bind(students); err != nil {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Bad Students Json"}))))
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
	return c.String(http.StatusOK, string(remErr(json.Marshal(r))))
}

func (h Controller) GetStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Invalid ID structure. ID should be positive integer"}))))
	}
	res, err := h.sh.GetStudent(uint(id))
	if errors.Is(err, common.StudentNotFoundError) || errors.Is(err, common.StudentClassError) {
		return c.String(http.StatusNotFound, string(remErr(json.Marshal(common.Error{Note: "Couldn't find requested student"}))))
	}
	return c.String(http.StatusOK, string(remErr(json.Marshal(res))))
}

func (h Controller) UpdateStudentScore(c echo.Context) error {
	data := new(model.StudentScore)
	if err := c.Bind(data); err != nil || data.ID == 0 {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Bad Json. ID should be positive integer"}))))
	}

	if err := h.sh.UpdateStudentScore(data.ID, data.Score); err != nil {
		if err == common.InvalidScoreError {
			return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Invalid score value. score should be between 0 an 20"}))))
		}
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Student not Found"}))))
	}
	return c.String(http.StatusOK, string(remErr(json.Marshal(common.Success{Note: "Student score updated"}))))
}

func (h Controller) UpdateStudentEmail(c echo.Context) error {
	data := new(model.StudentEmail)
	if err := c.Bind(data); err != nil || data.ID == 0 {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Bad Json. ID should be positive integer"}))))
	}

	if err := h.sh.UpdateStudentEmail(data.ID, data.Email); err != nil {
		if err == common.InvalidEmailError {
			return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Invalid email address"}))))
		}
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Student not Found"}))))
	}
	return c.String(http.StatusOK, string(remErr(json.Marshal(common.Success{Note: "Student email updated"}))))
}

func (h Controller) DeleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Invalid ID structure. ID should be positive integer"}))))
	}
	if err := h.sh.DeleteStudent(uint(id)); errors.Is(err, common.StudentNotFoundError) {
		return c.String(http.StatusNotFound, string(remErr(json.Marshal(common.Error{Note: "Couldn't find requested student"}))))
	}
	return c.String(http.StatusOK, "")
}
