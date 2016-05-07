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

func GetSubjects(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	showSubjects(w, r, "")
}

func NewSubject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	template, err := template.ParseFiles("templates/base.html", "templates/subject/new.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func CreateSubject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	// Parse and validate.
	short, name, splitClass := parseSubjectData(r)
	valid, message := validateSubjectData(short, name, splitClass)
	if valid {
		subject := model.Subject{Short: short}
		if subject.Exists() {
			message = fmt.Sprintf("Es gibt bereits ein Fach mit dem Kürzel %s.", short)
			valid = false
		}
	}

	// Render message if the data is not valid.
	if !valid {
		template, err := template.ParseFiles("templates/base.html", "templates/subject/new.html")
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		err = template.Execute(w, simpleMessage(message, false))
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		return
	}

	// Create the subject.
	subject := model.Subject{Short: short, Name: name, SplitClass: splitClass}
	subject.Create()

	// Render all subjects
	message = fmt.Sprintf("%s wurde gespeichert.", name)
	showSubjects(w, r, message)
}

func showSubjects(w http.ResponseWriter, r *http.Request, message string) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	// Prepare the template data with all subjects.
	templateData := struct {
		generalTemplateData
		Subjects []model.Subject
	}{Subjects: model.ReadAllSubjects()}
	// Add the message if there is one.
	if message != "" {
		templateData.Messages = []templateMessage{templateMessage{Text: message, Positive: true}}
	}

	// Render.
	template, err := template.ParseFiles("templates/base.html", "templates/subject/index.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	err = template.Execute(w, &templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func GetSubject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	subject := model.Subject{Short: short}
	if !subject.Exists() {
		http.NotFound(w, r)
		return
	}

	// Execute template with subject as template data.
	template, err := template.ParseFiles("templates/base.html", "templates/subject/get.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	templateData := struct {
		generalTemplateData
		model.Subject
	}{Subject: subject}

	err = template.Execute(w, &templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func EditSubject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	subject := model.Subject{Short: short}
	if !subject.Exists() {
		http.NotFound(w, r)
		return
	}

	// Execute template with subject as template data.
	template, err := template.ParseFiles("templates/base.html", "templates/subject/edit.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	templateData := struct {
		generalTemplateData
		model.Subject
	}{}
	templateData.Subject = subject

	err = template.Execute(w, &templateData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func UpdateSubject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	subject := model.Subject{Short: short}
	if !subject.Exists() {
		http.NotFound(w, r)
		return
	}

	// Parse and validate subject data.
	nshort, name, splitClass := parseSubjectData(r)
	valid, message := validateSubjectData(short, name, splitClass)
	if !valid {
		http.Error(w, message, http.StatusNotAcceptable)
		return
	}

	// Update subject and send the new URL back to the client.
	// It might have changed with an update of short.
	updSubject := model.Subject{Short: nshort, Name: name, SplitClass: splitClass}
	updSubject.UpdateShort(short)
	w.Write([]byte(fmt.Sprintf("/subject/%s", nshort)))
}

func DeleteSubject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirected, _ := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	short := html.EscapeString(params.ByName("short"))
	subject := model.Subject{Short: short}
	if subject.Exists() {
		subject.Delete()
	} else {
		http.NotFound(w, r)
		return
	}
}

func parseSubjectData(r *http.Request) (short, name string, splitClass bool) {
	r.ParseForm()
	short = html.EscapeString(r.Form.Get("short"))
	name = html.EscapeString(r.Form.Get("name"))
	_splitClass := html.EscapeString(r.Form.Get("splitClass"))
	if _splitClass == "true" {
		splitClass = true
	}
	return
}

func validateSubjectData(short, name string, splitClass bool) (valid bool, message string) {
	if len(short) == 0 {
		return false, "Das Kürzel ist zu kurz."
	} else if len(name) == 0 {
		return false, "Der Name ist zu kurz."
	}
	return true, ""
}
