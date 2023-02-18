package main

/*
NOTES: use html/template to pass variables to html: https://stackoverflow.com/questions/27971240/how-to-pass-just-a-variable-not-a-struct-member-into-text-html-template-golan

	- start docker service with `sudo service docker start`
*/

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var working_dir, static_dir, port, local, session_key string
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var controlPassword string

type User struct {
	Name string
}

type Answer struct {
	Id     int
	Body   string
	Votes  int
	Author string
}

type Question struct {
	Id   int
	Body string
}

type PageData struct {
	Question       Question
	RemainingVotes int
	Answered       bool
}

type QuestionRecord struct {
	Id            int
	Body          string
	Author        string
	DateSubmitted string
}

type VoteResponse struct {
	Result         string
	Votes          int
	RemainingVotes int
}

var DB *sql.DB

func main() {
	// Load .env file
	godotenv.Load()
	local = os.Getenv("RUN-LOCAL") // set to true if running on local machine
	session_key = os.Getenv("SESSION-KEY")
	controlPassword = os.Getenv("CONTROL-PASSWORD")

	if controlPassword == "" {
		panic("must set control password")
	}

	var srv *http.Server
	r := mux.NewRouter()

	if local == "false" {
		port = "80"
		working_dir = "/usr/"
		srv = &http.Server{
			Handler: r,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	} else {
		port = "9990"
		working_dir = "../"
		srv = &http.Server{
			Handler: r,
			Addr:    "localhost:9990",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	}

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(working_dir)))

	// Define routes
	r.HandleFunc("/", home)
	r.HandleFunc("/control", control)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/api/answers", get_answers).Methods("GET")
	r.HandleFunc("/api/answers", post_answer).Methods("POST")
	r.HandleFunc("/api/approve", approve_answer).Methods("POST")
	r.HandleFunc("/api/vote", vote).Methods("POST")
	r.HandleFunc("/api/questions", get_questions).Methods("GET")
	r.HandleFunc("/api/questions", post_question).Methods("POST")
	r.HandleFunc("/loaderio-d76bfe3fee5c082595ab976a8b88ed42/", loader_auth)

	//connect to db
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", "qotd", os.Getenv("DB-PASSWORD"), "qotd")
	var err error
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	log.Fatal(srv.ListenAndServe())
	defer DB.Close()
}
