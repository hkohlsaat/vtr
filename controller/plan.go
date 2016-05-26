package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/hkohlsaat/vtr/model"
	"github.com/julienschmidt/httprouter"
)

func PostPlan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(65536)
	if r.Form.Get("passwort") != os.Args[2] {
		http.Error(w, "falsches Passwort", http.StatusUnauthorized)
		return
	}
	file, _, err := r.FormFile("vertretungsplan")
	if err != nil {
		log.Printf("error recieving file from upload: %v\n", err)
		return
	}
	go processPlan(file)
}

func processPlan(file multipart.File) {
	defer file.Close()

	plan, err := model.ToPlan(file)
	if err != nil {
		log.Printf("can't make an object of the plan: %v\n", err)
		return
	}

	file.Seek(0, 0)
	bytes, _ := ioutil.ReadAll(file)
	plan.Create(bytes)

	for _, part := range plan.Parts {
		for _, s := range part.Substitutions {
			if s.SubstTeacher.Name == "" && s.SubstTeacher.Short != "" {
				unknownTeacher := model.UnknownTeacher{Short: s.SubstTeacher.Short}
				unknownTeacher.Create()
			}
			if s.InstdTeacher.Name == "" && s.InstdTeacher.Short != "" {
				unknownTeacher := model.UnknownTeacher{Short: s.InstdTeacher.Short}
				unknownTeacher.Create()
			}
			if s.InstdSubject.Name == "" && s.InstdSubject.Short != "" {
				unknownSubject := model.UnknownSubject{Short: s.InstdSubject.Short}
				unknownSubject.Create()
			}
		}
	}
}

func GetPlan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lastPlan := model.LastPlanJSON()
	if lastPlan == "" {
		http.NotFound(w, r)
		return
	} else {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		fmt.Fprintln(w, lastPlan)
	}
}
