package routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ie-backend-project/common"
	"ie-backend-project/controller"
	"ie-backend-project/handler"
	"strconv"
)

func NewRouter(port uint, basePath string, courseHandler *handler.CourseHandler, studentHandler *handler.StudentHandler) error {
	if port < 1000 || port > 65535 {
		return errors.New("unacceptable port num")
	}

	e := echo.New()
	c := controller.NewController(courseHandler, studentHandler)

	e.POST(basePath+"/register", c.Register)
	e.POST(basePath+"/login", c.Login)
	e.POST(basePath+"/logout", c.Logout)

	cg := e.Group(basePath + "/course")
	sg := e.Group(basePath + "/student")
	config := middleware.JWTConfig{
		Claims:     &common.JWTCustomClaims{},
		SigningKey: common.JWTKey,
	}
	cg.Use(middleware.JWTWithConfig(config))
	sg.Use(middleware.JWTWithConfig(config))

	cg.POST("/new", c.NewCourse)
	cg.POST("/delete", c.DeleteCourse)
	cg.GET("/:id", c.GetCourse)
	sg.POST("/new", c.NewStudent)
	sg.POST("/delete", c.DeleteStudent)
	sg.POST("/update/score", c.UpdateStudentScore)
	sg.POST("/update/email", c.UpdateStudentEmail)
	sg.GET("/:id", c.GetStudent)

	return e.Start(":" + strconv.Itoa(int(port)))
}
