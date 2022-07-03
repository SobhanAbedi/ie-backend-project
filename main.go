package main

import (
	"ie-backend-project/handler"
	"ie-backend-project/mailer"
	"ie-backend-project/routes"
)

const (
	mailUser = "abedi.sobhan2000@gmail.com"
	mailPass = "xxxxxxxxxxx"
)

func safe[T any](output T, err error) T {
	if err != nil {
		panic(err)
	}
	return output
}

func main() {
	courseHandler := safe(handler.NewCourseHandler("main.db"))
	studentHandler := safe(handler.NewStudentHandler("main.db"))
	studentMailer := mailer.NewMailer(mailUser, mailPass)
	safe(0, routes.NewRouter(8080, "/api", courseHandler, studentHandler, studentMailer))
}
