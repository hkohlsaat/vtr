package main

import (
	"net/http"
	"os"

	"github.com/hkohlsaat/vtr/controller"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", controller.Index)
	router.GET("/signup", controller.GetSignup)
	router.POST("/signup", controller.PostSignup)
	router.GET("/login", controller.GetLogin)
	router.POST("/login", controller.PostLogin)

	router.GET("/teachers", controller.GetTeachers)
	router.GET("/teachers/new", controller.NewTeacher)
	router.POST("/teachers", controller.CreateTeacher)
	router.GET("/teachers/upload", controller.NewTeachers)
	router.POST("/teachers/upload", controller.CreateTeachers)
	router.GET("/teacher/:short", controller.GetTeacher)
	router.GET("/teacher/:short/edit", controller.EditTeacher)
	router.PUT("/teacher/:short", controller.UpdateTeacher)
	router.DELETE("/teacher/:short", controller.DeleteTeacher)

	router.GET("/subjects", controller.GetSubjects)
	router.GET("/subjects/new", controller.NewSubject)
	router.POST("/subjects", controller.CreateSubject)
	router.GET("/subjects/upload", controller.NewSubjects)
	router.POST("/subjects/upload", controller.CreateSubjects)
	router.GET("/subject/:short", controller.GetSubject)
	router.GET("/subject/:short/edit", controller.EditSubject)
	router.PUT("/subject/:short", controller.UpdateSubject)
	router.DELETE("/subject/:short", controller.DeleteSubject)

	router.GET("/plan", controller.GetPlan)
	router.POST("/plan", controller.PostPlan)

	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	port := os.Args[1]
	http.ListenAndServe(":"+port, router)
}
