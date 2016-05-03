package main

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/hkohlsaat/vtr/model"

	"github.com/julienschmidt/httprouter"
)

type generalTemplateData struct {
	Messages []templateMessage
}

type templateMessage struct {
	Text     string
	Positive bool
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/signup", GetSignup)
	router.POST("/signup", PostSignup)
	router.GET("/login", GetLogin)
	router.POST("/login", PostLogin)
	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	http.ListenAndServe(":9000", router)
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirected, session := ensureLoggedIn(w, r)
	if redirected {
		return
	}

	template, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	err = template.Execute(w, simpleMessage(fmt.Sprintf("%+v", session), true))
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func ensureLoggedIn(w http.ResponseWriter, r *http.Request) (redirected bool, session model.Session) {
	if model.CountUsers() == 0 {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return true, model.Session{}
	}

	cookie, err := r.Cookie(model.SStore.CookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return true, model.Session{}
	}

	sid, _ := url.QueryUnescape(cookie.Value)
	session, ok := model.SStore.Session(sid)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return true, model.Session{}
	}

	return false, session
}

func GetSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if model.CountUsers() > 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	template, err := template.ParseFiles("templates/base.html", "templates/signup.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func PostSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if model.CountUsers() > 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get information from form.
	r.ParseForm()
	var username = html.EscapeString(r.Form.Get("username"))
	var password = r.Form.Get("password")

	var validData bool
	var message string
	if model.UsernameTaken(username) {
		message = "Dieser Nutzername ist bereits vergeben."
	} else if len(password) < 3 {
		message = "Das Passwort ist zu kurz."
	} else {
		validData = true
	}
	if !validData {
		template, err := template.ParseFiles("templates/base.html", "templates/signup.html")
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		err = template.Execute(w, simpleMessage(message, false))
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		return
	}

	// Create new user
	user := model.User{Name: username}
	user.Create(password)
	// Login the user
	login(w, user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	template, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get information from form.
	r.ParseForm()
	var username = html.EscapeString(r.Form.Get("username"))
	var password = r.Form.Get("password")
	// Authenticate the user.
	var user = model.User{Name: username}
	var authError = user.GetWithPassword(password)
	if authError != nil {
		// Authetification went wrong.
		// Display a message to the user.
		var message string
		switch {
		case authError == model.ErrNoMatchNamePassword:
			message = "Das Passwort passt leider nicht zum Nutzernamen."
		case authError == model.ErrNoSuchUser:
			message = "Diesen Nutzernamen gibt es leider nicht."
		default:
			message = "Etwas Unvorhergesehenes ist passiert. Bitte versuche es noch einmal."
		}

		// Execute the template.
		template, err := template.ParseFiles("templates/base.html", "templates/login.html")
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		err = template.Execute(w, simpleMessage(message, false))
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	} else {
		login(w, user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func login(w http.ResponseWriter, user model.User) {
	// Create the new session.
	sess := model.SStore.NewSession(user)
	// Prepare the session cookie.
	cName := model.SStore.CookieName
	cValue := url.QueryEscape(sess.Id)
	cMaxAge := model.SStore.MaxAge
	cookie := http.Cookie{Name: cName, Value: cValue, Path: "/", HttpOnly: true, MaxAge: cMaxAge}
	// Set the cookie and redirect.
	http.SetCookie(w, &cookie)
}

func simpleMessage(message string, positive bool) *generalTemplateData {
	templateData := generalTemplateData{Messages: []templateMessage{templateMessage{Text: message, Positive: positive}}}
	return &templateData
}
