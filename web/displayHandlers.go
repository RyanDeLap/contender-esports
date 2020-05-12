package web

import (
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session := NewSession(r, w)
	loggedIn, _ := session.Get("logged_in")
	if loggedIn != nil && loggedIn.(bool) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}

	err := HtmlTemplate.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		panic(err)
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	session := NewSession(r, w)
	loggedIn, _ := session.Get("logged_in")
	if loggedIn == nil || !loggedIn.(bool) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := struct {
		LoggedIn bool
		Username string
		Email    string
	}{}

	err := HtmlTemplate.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		//log.Println(data)
		log.Println(err)
	}
}
