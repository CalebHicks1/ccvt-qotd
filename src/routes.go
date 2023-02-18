package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, working_dir+"static/index.html")

	// Get a cokkie session: https://github.com/gorilla/sessions
	// session, _ := store.Get(r, "session")

	// Build template from index.html
	templ, err := template.ParseFiles("../static/templates/index.html")
	if err != nil {
		fmt.Print(err)
	}

	// usernameString, ok := session.Values["userName"].(string)

	// get most recent question
	result := DB.QueryRow("select id, body from questions order by date_submitted desc limit 1")

	var resultQuestion QuestionRecord
	err = result.Scan(&resultQuestion.Id, &resultQuestion.Body)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}

	question := Question{
		Id:   resultQuestion.Id,
		Body: resultQuestion.Body,
	}

	data := PageData{
		Question:       question,
		RemainingVotes: 3,
	}
	// if ok {
	// 	data.Name = usernameString
	// }

	templ.Execute(w, data)
}

// API Routes

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

func get_answers(w http.ResponseWriter, r *http.Request) {

	a := Answer{1, "test answer", 10, "server"}
	b, err := json.Marshal(a)

	if err != nil {

		fmt.Fprintf(w, "json error")
	}
	fmt.Fprintf(w, "[%s]", string(b))
}

func get_questions(w http.ResponseWriter, r *http.Request) {

	result := DB.QueryRow("select id, body from questions order by date_submitted desc limit 1")

	var resultQuestion QuestionRecord
	err := result.Scan(&resultQuestion.Id, &resultQuestion.Body)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	q := Question{
		Id:   resultQuestion.Id,
		Body: resultQuestion.Body,
	}
	b, err := json.Marshal(q)
	fmt.Fprintf(w, "%s", string(b))
}
