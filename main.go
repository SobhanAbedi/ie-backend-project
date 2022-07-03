package main

import (
	"ie-backend-project/handler"
	"ie-backend-project/routes"
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
	//ieCourse := safe(courseHandler.NewCourse("Internet Engineering", "Parham Alvani"))
	//fmt.Println(ieCourse)
	//osCourse := safe(courseHandler.NewCourse("Operating Systems", "Ahmad Seyed Javadi"))
	//fmt.Println(osCourse)
	//std1 := safe(studentHandler.NewStudent("Sobhan", "Abedi", "abdi.sobhan2000@gmail.com", 18, ieCourse))
	//fmt.Println(std1)
	//safe(0, studentHandler.UpdateStudentScore(std1.ID, 0))
	//safe(0, studentHandler.UpdateStudentEmail(std1.ID, "abedi.sobhan@aut.ac.ir"))
	//updatedStd := safe(studentHandler.GetStudent(std1.ID))
	//fmt.Println(updatedStd)
	//studentHandler.DeleteStudent(updatedStd.ID)
	//courseHandler.DeleteCourse(ieCourse.ID)
	//courseHandler.DeleteCourse(osCourse.ID)
	safe(0, routes.NewRouter(8080, "/api", courseHandler, studentHandler))
}
