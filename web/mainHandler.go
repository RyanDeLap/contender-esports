package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//func SendDiscordMessage(w http.ResponseWriter, r *http.Request) {
//	var data struct {
//		Message string `json:"message"`
//	}
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&data); err != nil {
//		RespondWithError(w, http.StatusOK, "Invalid request payload")
//		return
//	}
//	defer r.Body.Close()
//
//	DiscordChannel <- data.Message
//}

func CreateCenterHandler(w http.ResponseWriter, r *http.Request) {

	//session := NewSession(r, w)
	//loggedIn, _ := session.Get("logged_in")
	//if loggedIn == nil || !loggedIn.(bool) {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	returnc
	//}

	var data struct {
		CenterName string `json:"centerName"`
		GuildID    string `json:"guildID"`
		ChannelID  string `json:"channelID"`
		GGLink     string `json:"gglink"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		RespondWithError(w, http.StatusOK, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	success := createCenter(data.CenterName, data.GuildID, data.ChannelID, data.GGLink)
	if !success {
		RespondWithError(w, http.StatusOK, "Error creating center.")
		return
	}

	RespondWithError(w, http.StatusOK, "")
	return

}

func UpdateCenterHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	centerIDRaw := vars["centerID"]
	centerID, _ := strconv.Atoi(centerIDRaw)

	var data struct {
		CenterName string `json:"centerName"`
		GuildID    string `json:"discordChannelID"`
		ChannelID  string `json:"discordGuildID"`
		GGLink     string `json:"ggLeapLink"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		RespondWithError(w, http.StatusOK, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	success := updateCenter(data.CenterName, data.GuildID, data.ChannelID, data.GGLink, centerID)
	if !success {
		RespondWithError(w, http.StatusOK, "Error creating center.")
		return
	}

	RespondWithError(w, http.StatusOK, "")
	return
}

func GetCenters(w http.ResponseWriter, r *http.Request) {
	//session := NewSession(r, w)
	//loggedIn, _ := session.Get("logged_in")
	//if loggedIn == nil || !loggedIn.(bool) {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}
	centers := getAllCenters()
	//pagesJson, _ := json.Marshal(centrs)

	RespondWithJson(w, http.StatusOK, centers)
	return
}

func InsertScheduleHandler(w http.ResponseWriter, r *http.Request) {
	//session := NewSession(r, w)
	//loggedIn, _ := session.Get("logged_in")
	//if loggedIn == nil || !loggedIn.(bool) {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}

	var data struct {
		TimeToPost string `json:"timeToPost"`
		DayOfWeek  string `json:"dayOfWeek"`
		Game       string `json:"game"`
		CenterID   int    `json:"centerID"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		RespondWithError(w, http.StatusOK, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	success := insertNewSchedule(data.TimeToPost, data.DayOfWeek, data.Game, data.CenterID)
	if !success {
		RespondWithError(w, http.StatusOK, "Error creating center.")
		return
	}

	clearScheduler()
	time.Sleep(1 * time.Second)
	setupScheduler() //Reload scheduler config.

	RespondWithError(w, http.StatusOK, "")
	return
}

func UpdateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	//session := NewSession(r, w)
	//loggedIn, _ := session.Get("logged_in")
	//if loggedIn == nil || !loggedIn.(bool) {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}

	var data struct {
		TimeToPost string `json:"timeToPost"`
		DayOfWeek  string `json:"dayOfWeek"`
		Game       string `json:"game"`
		ScheduleID int    `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		RespondWithError(w, http.StatusOK, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	success := updateSchedule(data.TimeToPost, data.DayOfWeek, data.Game, data.ScheduleID)
	if !success {
		RespondWithError(w, http.StatusOK, "Error creating center.")
		return
	}

	clearScheduler()
	time.Sleep(1 * time.Second)
	setupScheduler() //Reload scheduler config.

	RespondWithError(w, http.StatusOK, "")
	return
}

type NewScheduleRequest struct {
	TimeToPost string `json:"time"`
	DayOfWeek  string `json:"day"`
	Game       string `json:"game"`
}

func SaveSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	//session := NewSession(r, w)
	//loggedIn, _ := session.Get("logged_in")
	//if loggedIn == nil || !loggedIn.(bool) {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}

	vars := mux.Vars(r)
	centerIDRaw := vars["centerID"]
	centerID, _ := strconv.Atoi(centerIDRaw)

	var data []NewScheduleRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
		RespondWithError(w, http.StatusOK, "Invalid request payload")
		return
	}

	success := deleteSchedules(centerID)
	if success == false {
		fmt.Println("Error deleting")
		RespondWithError(w, http.StatusOK, "Error deleting.")
		return
	}

	for _, schedule := range data {
		success := insertNewSchedule(schedule.TimeToPost, schedule.DayOfWeek, schedule.Game, centerID)
		if success == false {
			RespondWithError(w, http.StatusOK, "Error inserting.")
			return
		}
	}

	defer r.Body.Close()

	clearScheduler()
	time.Sleep(1 * time.Second)
	setupScheduler() //Reload scheduler config.

	RespondWithError(w, http.StatusOK, "")
	return
}
