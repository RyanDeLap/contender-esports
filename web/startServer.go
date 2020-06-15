package web

import (
	"contender/discord"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tdewolff/minify"
)

type PanelConfig struct {
	MysqlUsername string `json:"mysql_username"`
	MysqlPassword string `json:"mysql_password"`
	MysqlHost     string `json:"mysql_host"`
	MysqlDb       string `json:"mysql_db"`
	DiscordApiKey string `json:"discord_api_key"`
}

var mysqlConn *sql.DB
var HtmlTemplate *template.Template
var DEBUG = true
var DiscordChannel chan discord.ChannelMessageInfo = make(chan discord.ChannelMessageInfo)
var oauthConfig oauth1.Config

func Start() {
	fmt.Println("Starting Contender eSports Web Panel")

	config := loadConfig()
	connectMySQL(config)

	if DEBUG {
		fmt.Println("Running in DEBUG mode!")
		if !userExists("root") {
			createUser("ryandelap", "asdfasdf44", "contact@ryandelap.me")
			fmt.Println("Created test user.")
		}
	}

	fmt.Println("Starting discord bot..." + config.DiscordApiKey)

	go discord.RunDiscordBot(DiscordChannel, config.DiscordApiKey)

	oauthConfig = oauth1.Config{
		ConsumerKey:    "",
		ConsumerSecret: "",
		CallbackURL:    "http://localhost:8080/oauth/twitter/callback",
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	fmt.Println("Loading schedulers")

	gob.Register(StuffToSave{})

	fmt.Println("Running server...")

	setupRouter() //Must be called last!
}

func loadConfig() PanelConfig {

	var config PanelConfig
	configFile, _ := os.Open("config.json")
	byteValue, _ := ioutil.ReadAll(configFile)

	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Println("Error while loading config.json:", err)
	}

	configFile.Close()

	return config
}

func connectMySQL(config PanelConfig) {
	var err error
	if mysqlConn, err = sql.Open("mysql", config.MysqlUsername+":"+config.MysqlPassword+"@("+config.MysqlHost+")/"+config.MysqlDb); err != nil {
		log.Fatal(err)
	}

	if err := mysqlConn.Ping(); err != nil {
		log.Fatal(err)
	}
}

func setupRouter() {
	var err error

	r := mux.NewRouter()
	m := minify.New()

	dynamicRouter := r.PathPrefix("/").Subrouter()
	//dynamicRouter.HandleFunc("/", endpoints.homeHandler)
	dynamicRouter.HandleFunc("/login", LoginHandler).Methods("GET")
	dynamicRouter.HandleFunc("/login", LoginHandlerPOST).Methods("POST")
	dynamicRouter.HandleFunc("/dashboard", DashboardHandler)

	//Business Logic
	//dynamicRouter.HandleFunc("/send-discord-message", SendDiscordMessage).Methods("POST")
	dynamicRouter.HandleFunc("/create-center", CreateCenterHandler).Methods("POST")

	dynamicRouter.HandleFunc("/insert-schedule", InsertScheduleHandler).Methods("POST")
	dynamicRouter.HandleFunc("/update-schedule", UpdateScheduleHandler).Methods("POST")
	dynamicRouter.HandleFunc("/save-schedules/{centerID}", SaveSchedulesHandler).Methods("POST")
	dynamicRouter.HandleFunc("/save-config/{centerID}", UpdateCenterHandler).Methods("POST")

	dynamicRouter.HandleFunc("/centers", GetCenters).Methods("GET")

	dynamicRouter.HandleFunc("/oauth/twitter/callback", twitterCallbackHandler)
	dynamicRouter.HandleFunc("/login_twitter", loginHandler).Methods("POST")

	dynamicRouter.HandleFunc("/save_facebook_page", saveFacebookPageHandler).Methods("POST")
	dynamicRouter.HandleFunc("/upconvert_token", upconvertUserTokenToExtendedTokenHandler).Methods("POST")

	r.PathPrefix("/assets/").Handler(
		m.Middleware(
			http.StripPrefix("/assets/",
				http.FileServer(http.Dir("static/assets/")),
			),
		),
	)

	recompileHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			HtmlTemplate, err = template.ParseGlob("static/*.html")
			if err != nil {
				log.Println("Failed loading html templates")
			}
			next.ServeHTTP(w, r)
		})
	}

	accessLogFile, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logWriter := accessLogFile

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CombinedLoggingHandler(logWriter, recompileHandler(r)),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
