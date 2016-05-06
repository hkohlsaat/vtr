package main

import (
	"net/http"

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
	router.GET("/teacher/:short", controller.GetTeacher)
	router.GET("/teacher/:short/edit", controller.EditTeacher)
	router.PUT("/teacher/:short", controller.UpdateTeacher)
	router.DELETE("/teacher/:short", controller.DeleteTeacher)

	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	http.ListenAndServe(":9000", router)
}
