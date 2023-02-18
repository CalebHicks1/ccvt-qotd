package main

/*
NOTES: use html/template to pass variables to html: https://stackoverflow.com/questions/27971240/how-to-pass-just-a-variable-not-a-struct-member-into-text-html-template-golan

	- start docker service with `sudo service docker start`
*/

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var working_dir, static_dir, port, local, session_key string
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type User struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, working_dir+"static/index.html")

	// Get a cokkie session: https://github.com/gorilla/sessions
	session, _ := store.Get(r, "session")

	// Build template from index.html
	templ, err := template.ParseFiles("../static/index.html")
	if err != nil {
		fmt.Print(err)
	}

	usernameString, ok := session.Values["userName"].(string)
	data := User{
		Name: usernameString,
	}
	if ok {
		data.Name = usernameString
	}

	templ.Execute(w, data)
}

func api_get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "running locally: %s", local)
}

func api_post(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("POST request from %s: username=%s\n", r.RemoteAddr, r.PostFormValue("username"))
	session, _ := store.Get(r, "session")
	session.Values["userName"] = r.PostFormValue("username")
	fmt.Printf(r.PostFormValue("username"))
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "https://calebhicks.net", http.StatusSeeOther)
}

func main() {
	// Load .env file
	godotenv.Load()
	local = os.Getenv("RUN-LOCAL") // set to true if running on local machine
	session_key = os.Getenv("SESSION-KEY")
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
	r.HandleFunc("/api", api_get).
		Methods("GET")
	r.HandleFunc("/api", api_post).
		Methods("POST")

	// Start the server

	log.Fatal(srv.ListenAndServe())
}
