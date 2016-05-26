package controller

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hkohlsaat/vtr/model"
	"github.com/julienschmidt/httprouter"
)

// GetTeachers serves the list of all teachers.
func GetTeachers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	showTeachers(w, r, "")
}

// NewTeacher serves the form to create a new teacher.
func NewTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(r.URL.Query().Get("short"))
	templateData := struct {
		generalTemplateData
		Short string
	}{Short: short}

	template, err := template.ParseFiles("templates/base.html", "templates/teacher/new.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	err = template.Execute(w, templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

// CreateTeacher creates a new teacher and serves the list of all teachers.
func CreateTeacher(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	// Parse and validate.
	r.ParseForm()
	short := html.EscapeString(r.Form.Get("short"))
	name := html.EscapeString(r.Form.Get("name"))
	sex := html.EscapeString(r.Form.Get("sex"))

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

	unknown := model.UnknownTeacher{Short: short}
	unknown.Delete()

	// Render all teachers
	message = fmt.Sprintf("%s wurde gespeichert.", short)
	showTeachers(w, r, message)
}

// showTeachers is a helper function to show a list of all teachers.
func showTeachers(w http.ResponseWriter, r *http.Request, message string) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	// Prepare the template data with all teachers.
	templateData := struct {
		generalTemplateData
		Teachers []model.Teacher
		Unknown  []model.UnknownTeacher
	}{Teachers: model.ReadAllTeachers(),
		Unknown: model.ReadAllUnknownTeachers()}
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

// NewTeachers serves the upload form to submit multiple teacher records.
func NewTeachers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	template, err := template.ParseFiles("templates/base.html", "templates/teacher/upload.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

// CreateTeachers creates multiple teachers from json upload and serves the list
// of all teachers.
func CreateTeachers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	r.ParseMultipartForm(16384)
	file, _, err := r.FormFile("teacherjson")
	if err != nil {
		log.Printf("error recieving teacher file from upload: %v\n", err)
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading teacher file: %v\n", err)
		return
	}

	type tmpTeacher struct {
		model.Teacher
		Compellation string
	}
	teachers := []tmpTeacher{}
	err = json.Unmarshal(bytes, &teachers)
	if err != nil {
		log.Printf("error unmarshaling teacher data: %v\n", err)
		return
	}
	for _, teacher := range teachers {
		if teacher.Compellation == "Herr" {
			teacher.Sex = "m"
		} else {
			teacher.Sex = "w"
		}
		teacher.Create()

		unknown := model.UnknownTeacher{Short: teacher.Short}
		unknown.Delete()
	}

	http.Redirect(w, r, "/teachers", http.StatusSeeOther)
}

// GetTeacher serves one teacher.
func GetTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if !teacher.Exists() {
		http.NotFound(w, r)
		return
	} else {
		teacher.Read()
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

// EditTeacher serves a form to edit a teacher.
func EditTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if !teacher.Exists() {
		http.NotFound(w, r)
		return
	}
	teacher.Read()

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

// UpdateTeacher updates a teacher with the uploaded information.
func UpdateTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if !teacher.Exists() {
		http.NotFound(w, r)
		return
	}

	// Parse and validate teacher data.
	r.ParseForm()
	nshort := html.EscapeString(r.Form.Get("short"))
	name := html.EscapeString(r.Form.Get("name"))
	sex := html.EscapeString(r.Form.Get("sex"))

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

// DeleteTeacher deletes a teacher and serves nothing (empty 200 OK response).
func DeleteTeacher(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	teacher := model.Teacher{Short: short}
	if teacher.Exists() {
		teacher.Delete()
	} else {
		http.NotFound(w, r)
		return
	}
}

func validateTeacherData(short, name, sex string) (valid bool, message string) {
	if sex != "m" && sex != "w" {
		return false, "Falsche Eingabe: Herr: \"m\", Frau: \"w\""
	} else if len(short) == 0 {
		return false, "Das Kürzel ist zu kurz."
	} else if len(name) == 0 {
		return false, "Der Name ist zu kurz."
	}
	return true, ""
}
