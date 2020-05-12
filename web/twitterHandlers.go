package web

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func twitterCallbackHandler(w http.ResponseWriter, req *http.Request) {
	session := NewSession(req, w)

	var requestSecretStr string
	if requestSecret, err := session.Get("twitterRequestSecret"); err == nil {
		requestSecretStr = requestSecret.(string)
	} else {
		log.Println("call back with missing request secret. err:", err)
		return
	}

	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(req)
	if err != nil {
		log.Println("Failed parsing auth callback:", err)
		return
	}

	accessToken, accessSecret, err := oauthConfig.AccessToken(requestToken, requestSecretStr, verifier)
	if err != nil {
		log.Println("Failed to get secret and token from response. err:", err)
		return
	}

	log.Println("got user access token:", accessToken, ", secret:", accessSecret)


	centerIDRaw, _ := session.Get("centerID")
	centerID, _ := strconv.Atoi(centerIDRaw.(string))
	fmt.Println("CenterID: ", centerID)



	success := updateTwitterOAUTH(accessToken, accessSecret, centerID)
	if !success {
		fmt.Println("Failed to insert OAUTH tokens")
	}

	http.Redirect(w, req, "/dashboard", http.StatusFound)
}


func loginHandler(w http.ResponseWriter, req *http.Request) {
	session := NewSession(req, w)


	requestToken, requestSecret, err := oauthConfig.RequestToken()
	if err != nil {
		log.Println("Error getting initial request token:", err)
		return
	}

	log.Println("request token:", requestToken, ", request secret:", requestSecret)

	authorizationURL, err := oauthConfig.AuthorizationURL(requestToken)
	if err != nil {
		log.Println("Failed to get authorization url:", err)
		return
	}

	err = req.ParseForm()
	if err != nil {
		panic(err)
	}

	session.Set("centerID", req.Form.Get("centerID"))


	session.Set("twitterRequestSecret", requestSecret)

	http.Redirect(w, req, authorizationURL.String(), http.StatusFound)
}

func sendTwitterTweeterTweet(centerInfo CenterInfo, status, imageFile string) error {
	api := anaconda.NewTwitterApiWithCredentials(centerInfo.TwitterAccessToken, centerInfo.TwitterSecret, oauthConfig.ConsumerKey, oauthConfig.ConsumerSecret)

	fileFd, err := os.Open(imageFile)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to open input file:", err))
	}
	fileContent, err := ioutil.ReadAll(fileFd)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to read in file contents:", err))
	}
	fileFd.Close()

	fileContentStr := base64.StdEncoding.EncodeToString(fileContent)

	media, err := api.UploadMedia(fileContentStr)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to upload image to twitter:", err))
	}

	log.Println("got twitter media id:", media.MediaID)

	twit, err := api.PostTweet(status, url.Values{
		"media_ids": []string{strconv.Itoa(int(media.MediaID))},
	})
	if err != nil {
		return errors.New(fmt.Sprint("Failed to post tweet to twitter:", err))
	}

	log.Println("tweet return:", twit)

	return nil

	//token := oauth1.NewToken(accessToken, accessSecret)
	//httpClient := oauthConfig.Client(oauth1.NoContext, token)
	//client := twitterClient.NewClient(httpClient)

	//tweet, resp, err := client.Statuses.Update("new tweet who dis", nil)
	//if err != nil {
	//	log.Println("Failed to send tweet. err:", err)
	//	return
	//}
}