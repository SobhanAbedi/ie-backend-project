package routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"ie-backend-project/controller"
	"ie-backend-project/handler"
	"net/http"
	"strconv"
)

type Router struct {
	port       uint
	bp         string //Base Path
	app        *echo.Echo
	controller *controller.Controller
}

func NewRouter(port uint, basePath string, courseHandler *handler.CourseHandler, studentHandler *handler.StudentHandler) (*Router, error) {
	if port < 1000 || port > 65535 {
		return nil, errors.New("unacceptable port num")
	}
	r := Router{port: port, bp: basePath, app: echo.New(), controller: controller.NewController(courseHandler, studentHandler)}

	r.app.GET(r.bp+"/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	r.app.POST(r.bp+"/register", r.controller.Register)
	r.app.POST(r.bp+"/login", r.controller.Login)
	r.app.POST(r.bp+"/logout", r.controller.Logout)
	r.app.POST(r.bp+"/new-course", r.controller.NewCourse)
	r.app.GET(r.bp+"/get-course/:id", r.controller.GetCourse)
	r.app.GET(r.bp+"/delete-course/:id", r.controller.DeleteCourse)
	//app.Get("/api/user", controllers.User)

	return &r, nil
}

func (r Router) Start() error {
	res := r.app.Start(":" + strconv.Itoa(int(r.port)))
	//r.app.Logger(res)
	return res
}