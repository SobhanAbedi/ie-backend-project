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
	res, err := h.ch.AddCourse(*course)
	if err != nil {
		if errors.Is(err, common.DuplicateCourseError) {
			return c.String(http.StatusForbidden, string(remErr(json.Marshal(common.Error{Note: "Course Already Exists"}))))
		}
		return c.String(http.StatusExpectationFailed, string(remErr(json.Marshal(common.Error{Note: "Couldn't add course"}))))
	}
	fmt.Println("Added", course)
	return c.String(http.StatusCreated, string(remErr(json.Marshal(common.Created{ID: res}))))
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, string(remErr(json.Marshal(common.Error{Note: "Invalid ID structure. ID should be positive integer"}))))
	}
	if err := h.ch.DeleteCourse(uint(id)); errors.Is(err, common.CourseNotFoundError) {
		return c.String(http.StatusNotFound, string(remErr(json.Marshal(common.Error{Note: "Couldn't find requested course"}))))
	}
	fmt.Println("Deleted Course", id)
	return c.String(http.StatusOK, "")
}
