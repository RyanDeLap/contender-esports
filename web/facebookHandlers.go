package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// FBAPI is the base URL for Facebook Graph API calls
	FBAPI = "https://graph.facebook.com"

	// FBVIDEOAPI is the base URL for uploading video via the Facebook Graph API
	FBVIDEOAPI = "https://graph-video.facebook.com"

	// FBAPIVERSION is the version of the Facebook API to use
	FBAPIVERSION = "v2.8"

	// FBAPPID is the APP ID from the Facebook API Dashboard
	FBAPPID = ""

	// FBAPPSECRET is the APP Secret from the Facebook API Dashboard
	FBAPPSECRET = ""

	// NoAttachment defines attachment type of nothing
	NoAttachment = 1 + iota
	// VideoAttachment defines attachment type of video
	VideoAttachment
	// PhotoAttachment defines attachment type of photo
	PhotoAttachment
	// LinkAttachment defines attachment type of link
	LinkAttachment
)

type Post struct {
	Title          string   `json:"title"`
	Message        string   `json:"message"`
	Attachment     string   `json:"attachment"`      // URL, depending on AttachmentType
	AttachmentType int      `json:"attachment_type"` // type of attachement
	Place          string   `json:"place"`           // Page ID of location associated with post (required to use tags)
	Tags           []string `json:"tags"`            // convert to comma separated values of user IDs ex: '1234,4566,6788'
}

type StuffToSave struct {
	FacebookPageAuthToken string `json:"facebook_page_auth_token"`
	FacebookPageID        string `json:"facebook_page_id"`
	FacebookPageName      string `json:"facebook_page_name"`
}

func saveFacebookPageHandler(w http.ResponseWriter, req *http.Request) {
	session := NewSession(req, w)
	var data = StuffToSave{}

	if dataIntf, err := session.Get("stuff_to_save"); err == nil && dataIntf != nil {
		data = dataIntf.(StuffToSave)
	} else {
		session.Set("stuff_to_save", data)
	}

	data.FacebookPageName = req.FormValue("page_name")
	data.FacebookPageID = req.FormValue("page_id")
	data.FacebookPageAuthToken = req.FormValue("page_token")

	centerIDRaw, _ := session.Get("centerID")
	centerID, _ := strconv.Atoi(centerIDRaw.(string))

	success := updateFacebookOAUTH(data.FacebookPageName, data.FacebookPageID, data.FacebookPageAuthToken, centerID)

	if !success {
		fmt.Println("error entering into database for facebook oauth toekns")
	}

	// save session again
	session.Set("stuff_to_save", data)
}

// idk if this is safe but /shrug
func upconvertUserTokenToExtendedTokenHandler(w http.ResponseWriter, req *http.Request) {
	tmpUserToken := req.FormValue("user_token")
	if tmpUserToken == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	endpoint := FBAPI + "/oauth/access_token?"
	params := []string{
		"grant_type=fb_exchange_token",
		"client_id=" + FBAPPID,
		"client_secret=" + FBAPPSECRET,
		"fb_exchange_token=" + tmpUserToken,
	}

	uri := endpoint + strings.Join(params, "&")

	client := http.Client{Timeout: time.Second * 5}
	body := &bytes.Buffer{}
	get, err := http.NewRequest(http.MethodGet, uri, body)
	if err != nil {
		log.Println("failed to create upgrade request:", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(get)
	if err != nil {
		log.Println("failed to get extended token from facebook:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("error reading response from facebook or bad status code:", resp.StatusCode, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var returnToken struct {
		AccessToken string `json:"access_token"`
	}

	err = json.Unmarshal(b, &returnToken)
	if err != nil {
		log.Println("error unmarshaling facebook extended token:", err, string(b))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("facebook return:", string(b))

	extendedToken := returnToken.AccessToken
	log.Println("upconverted token from:", tmpUserToken, ", to:", extendedToken)

	w.Header().Add("Content-Type", "application/json")

	// why marshal when u can just write manually :pepelaugh:
	w.Write([]byte(`{"token": "` + extendedToken + `"}`))
}

func postToFacebook(centerInfo CenterInfo, post Post) error {
	client := http.Client{Timeout: time.Second * 30}

	var edge string
	var api = FBAPI
	var params []string
	defaultParams := []string{
		"published=1",
		"access_token=" + centerInfo.FacebookPageAuthtoken,
	}

	switch post.AttachmentType {
	case VideoAttachment:
		params = []string{
			"title=" + url.QueryEscape(post.Title),
			"description=" + url.QueryEscape(post.Message),
			"file_url=" + post.Attachment,
		}

		edge = "videos"
		api = FBVIDEOAPI

	case PhotoAttachment:
		params = []string{
			"caption=" + url.QueryEscape(post.Message),
			"url=" + post.Attachment,
		}

		edge = "photos"

	case LinkAttachment:
		params = []string{
			"message=" + url.QueryEscape(post.Message),
			"link=" + post.Attachment,
		}

		edge = "feed"

	case NoAttachment:
		params = []string{
			"message=" + url.QueryEscape(post.Message),
		}

		edge = "feed"
	}

	uri := fmt.Sprintf("%s/v6.0/%s/%s?", api, centerInfo.FacebookPageID, edge)

	body := &bytes.Buffer{}
	params = append(defaultParams, params...)
	scheduleReq, err := http.NewRequest(http.MethodPost, uri+strings.Join(params, "&"), body)
	if err != nil {
		return errors.New(fmt.Sprint("failed creating facebook post request:", err))
	}

	log.Println("fb post url:", scheduleReq.URL)

	resp, err := client.Do(scheduleReq)
	if err != nil {
		return errors.New(fmt.Sprint("failed sending facebook post request:", err))
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprint("failed reading facebook post response:", err))
	}

	log.Println("post return:", string(b))

	var fbResponse map[string]interface{}
	err = json.Unmarshal(b, &fbResponse)
	if err != nil {
		return errors.New(fmt.Sprint("failed unmarshaling facebook post response:", err))
	}

	if fbResponse["error"] != nil {
		return errors.New(fmt.Sprint("facebook post response error:", fbResponse["error"]))
	}

	return nil
}
