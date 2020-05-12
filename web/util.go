package web

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
)

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isValidPassword(password string) bool {
	if len(password) >= 8 && len(password) <= 25 {
		return true
	}
	return false
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJson(w, code, map[string]string{"error": message})
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func isValidEmail(email string) bool {
	if len(email) > 45 {
		return false
	}

	return emailRegex.MatchString(email)
}

type ImgurResponse struct {
	Data struct {
		Link  string `json:"link"`
		Error string `json:"error"`
	} `json:"data"`
	Status  int  `json:"status"`
	Success bool `json:"success"`
}

func parseImgurResult(input string) *ImgurResponse {
	var res ImgurResponse
	err := json.Unmarshal([]byte(input), &res)
	if err != nil {
		log.Println("Failed parsing imgur result")
		return nil
	}
	return &res
}