package controller

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"

	"github.com/hkohlsaat/vtr/model"
	"github.com/julienschmidt/httprouter"
)

func GetTeachers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	showTeachers(w, r, "")
}

func NewTeacher(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	template, err := template.ParseFiles("templates/base.html", "templates/teacher/new.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func CreateTeacher(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse and validate.
	short, name, sex := parseTeacherData(r)
	valid, message := validateTeacherData(short, name, sex)
	if valid {
		teacher := model.Teacher{Short: short}
		if teacher.Exists() {
			message = fmt.Sprintf("Es gibt bereits einen Lehrer mit dem Kürzel %s.", short)
			valid = false
		}
	}

	// Render message if the data is not valid.
	if !valid {
		template, err := template.ParseFiles("templates/base.html", "templates/teacher/new.html")
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		err = template.Execute(w, simpleMessage(message, false))
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		return
	}

	// Create the teacher.
	teacher := model.Teacher{Short: short, Name: name, Sex: sex}
	teacher.Create()
	// Render all teachers
	message = fmt.Sprintf("%s wurde gespeichert.", short)
	showTeachers(w, r, message)
}

func showTeachers(w http.ResponseWriter, r *http.Request, message string) {
	// Prepare the template data with all teachers.
	templateData := struct {
		generalTemplateData
		Teachers []model.Teacher
	}{Teachers: model.ReadAllTeachers()}
	// Add the message if there is one.
	if message != "" {
		templateData.Messages = []templateMessage{templateMessage{Text: message, Positive: true}}
	}

	// Render.
	template, err := template.ParseFiles("templates/base.html", "templates/teacher/index.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	err = template.Execute(w, &templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func GetTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if !teacher.Exists() {
		http.NotFound(w, r)
		return
	}

	// Execute template with teacher as template data.
	template, err := template.ParseFiles("templates/base.html", "templates/teacher/get.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	templateData := struct {
		generalTemplateData
		model.Teacher
	}{Teacher: teacher}

	err = template.Execute(w, &templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func EditTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if !teacher.Exists() {
		http.NotFound(w, r)
		return
	}

	// Execute template with teacher as template data.
	template, err := template.ParseFiles("templates/base.html", "templates/teacher/edit.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	templateData := struct {
		generalTemplateData
		model.Teacher
	}{}
	templateData.Teacher = teacher

	err = template.Execute(w, &templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func UpdateTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if !teacher.Exists() {
		http.NotFound(w, r)
		return
	}

	// Parse and validate teacher data.
	nshort, name, sex := parseTeacherData(r)
	valid, message := validateTeacherData(short, name, sex)
	if !valid {
		http.Error(w, message, http.StatusNotAcceptable)
		return
	}

	// Update teacher and send the new URL back to the client.
	// It might have changed with an update of short.
	updTeacher := model.Teacher{Short: nshort, Name: name, Sex: sex}
	updTeacher.UpdateShort(short)
	w.Write([]byte(fmt.Sprintf("/teacher/%s", nshort)))
}

func DeleteTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if teacher.Exists() {
		teacher.Delete()
	} else {
		http.NotFound(w, r)
		return
	}
}

func parseTeacherData(r *http.Request) (short, name string, sex rune) {
	r.ParseForm()
	short = html.EscapeString(r.Form.Get("short"))
	name = html.EscapeString(r.Form.Get("name"))
	_sex := html.EscapeString(r.Form.Get("sex"))
	if len(_sex) > 0 {
		sex = rune(_sex[0])
	}
	return
}

func validateTeacherData(short, name string, sex rune) (valid bool, message string) {
	if sex != 'm' && sex != 'w' {
		return false, "Falsche Eingabe: Herr: 'm', Frau: 'w'"
	} else if len(short) == 0 {
		return false, "Das Kürzel ist zu kurz."
	} else if len(name) == 0 {
		return false, "Der Name ist zu kurz."
	}
	return true, ""
}
