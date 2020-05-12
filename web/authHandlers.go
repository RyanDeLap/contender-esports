package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func LoginHandlerPOST(w http.ResponseWriter, r *http.Request) {
	session := NewSession(r, w)
	loggedIn, _ := session.Get("logged_in")
	if loggedIn != nil && loggedIn.(bool) {
		RespondWithError(w, http.StatusOK, "Already logged in")
		return
	}

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginData); err != nil {
		RespondWithError(w, http.StatusOK, "Invalid request payload")
		fmt.Println("bad login")
		return
	}

	defer r.Body.Close()
	fmt.Println("%x", loginData)

	if strings.TrimSpace(loginData.Username) == "" {
		RespondWithError(w, http.StatusOK, "Invalid username or password.")
		return
	}

	if loginData.Password == "" {
		RespondWithError(w, http.StatusOK, "Invalid username or password.")
		return
	}

	if !CheckUserPassword(loginData.Username, loginData.Password) {
		RespondWithError(w, http.StatusOK, "Invalid username or password.")
		return
	}

	session.Set("logged_in", true)
	session.Set("username", loginData.Username)
	RespondWithError(w, http.StatusOK, "")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session := NewSession(r, w)
	loggedIn, _ := session.Get("logged_in")
	if loggedIn == nil || !loggedIn.(bool) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
