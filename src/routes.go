package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {

	// Build template from index.html
	templ, err := template.ParseFiles("../static/templates/index.html")
	if err != nil {
		fmt.Print(err)
	}

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
		RemainingVotes: get_remaining_votes(w, r),
		Answered:       get_question_answered(w, r),
	}

	templ.Execute(w, data)
}

func top(w http.ResponseWriter, r *http.Request) {

	// Build template from index.html
	templ, err := template.ParseFiles("../static/templates/top.html")
	if err != nil {
		fmt.Print(err)
	}

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
		Question: question,
	}

	templ.Execute(w, data)
}

func control(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, working_dir+"static/index.html")

	// Get a cookie session: https://github.com/gorilla/sessions
	session, _ := store.Get(r, "session")
	var logged_in bool
	if session.Values["logged_in"] != nil {

		logged_in, _ = session.Values["logged_in"].(bool)
	} else {
		logged_in = false
	}

	if logged_in {
		// Build template from index.html
		templ, err := template.ParseFiles("../static/templates/control.html")
		if err != nil {
			fmt.Print(err)
		}

		// get most recent question
		result := DB.QueryRow("select id, body from questions order by date_submitted desc limit 1")

		var resultQuestion QuestionRecord
		err = result.Scan(&resultQuestion.Id, &resultQuestion.Body)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}

		question := Question{
			Id:   int(resultQuestion.Id),
			Body: resultQuestion.Body,
		}

		data := PageData{
			Question:       question,
			RemainingVotes: 3,
		}

		templ.Execute(w, data)
	} else {
		templ, err := template.ParseFiles("../static/templates/login.html")
		if err != nil {
			fmt.Print(err)
		}

		templ.Execute(w, nil)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("password") == controlPassword {
		session, _ := store.Get(r, "session")
		session.Values["logged_in"] = true
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/control", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/control", http.StatusTemporaryRedirect)
	}
}

// Testing
func loader_auth(w http.ResponseWriter, r *http.Request) {
	token := "loaderio-d76bfe3fee5c082595ab976a8b88ed42"
	fmt.Fprintf(w, "%s", token)
}

func get_question_answered(w http.ResponseWriter, r *http.Request) bool {
	result := DB.QueryRow("select id from questions order by date_submitted desc limit 1")

	var question_id int
	err := result.Scan(&question_id)
	if err != nil {
		fmt.Print("error")
		return false
	}
	// fmt.Printf("current question id: %d", question_id)
	session, _ := store.Get(r, "session")
	latest_question_answered := session.Values["latest_question"]
	// if the user hasn't seen this question, return false
	if latest_question_answered == nil || session.Values["answered"] == nil || latest_question_answered.(int) < question_id {
		session.Values["latest_question"] = question_id
		session.Values["answered"] = false
		err := session.Save(r, w)
		if err != nil {
			fmt.Printf("error saving session")
		}
		fmt.Printf("false")
		return false
	} else {
		return session.Values["answered"].(bool)
	}

}

func get_remaining_votes(w http.ResponseWriter, r *http.Request) int {
	result := DB.QueryRow("select id from questions order by date_submitted desc limit 1")

	var question_id int
	err := result.Scan(&question_id)
	if err != nil {
		fmt.Print("error")
		return 0
	}
	// fmt.Printf("current question id: %d", question_id)
	session, _ := store.Get(r, "session")
	latest_question_answered := session.Values["latest_question"]
	if _, ok := latest_question_answered.(string); ok {
		latest_question_answered = 0
	}
	if latest_question_answered == nil || latest_question_answered.(int) < question_id {
		session.Values["latest_question"] = question_id
		session.Values["remaining_votes"] = 3
		session.Values["answered"] = false
		err := session.Save(r, w)
		if err != nil {
			fmt.Printf("error saving session")
		}
	}
	remaining_votes := session.Values["remaining_votes"]
	if remaining_votes != nil {
		return remaining_votes.(int)
	}
	session.Values["remaining_votes"] = 3
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("error saving session")
	}
	return 3
}

// API Routes

func get_answers(w http.ResponseWriter, r *http.Request) {
	var queryResult *sql.Rows
	var err error
	question_id := r.FormValue("question_id")
	if r.FormValue("source") == "control" {
		session, _ := store.Get(r, "session")
		if session.Values["logged_in"] != nil && session.Values["logged_in"].(bool) {
			queryResult, err = DB.Query("select id, body, author, votes from answers where question_id=? and approved=0", question_id)
		}
	} else if r.FormValue("source") == "home" {
		queryResult, err = DB.Query("select id, body, author, votes from answers where question_id=? and approved=1 order by date_submitted ASC", question_id)
	} else {
		queryResult, err = DB.Query("select id, body, author, votes from answers where question_id=? and approved=1 order by votes DESC", question_id)
	}
	if err != nil {
		fmt.Print(err.Error()) // proper error handling instead of panic in your app
	}
	if queryResult != nil {

		var answers []Answer
		for queryResult.Next() {
			var a Answer
			err = queryResult.Scan(
				&a.Id,
				&a.Body,
				&a.Author,
				&a.Votes,
			)
			answers = append(answers, a)
		}
		responseBytes, err := json.Marshal(answers)
		if err != nil {
			fmt.Printf("json error")
		}
		fmt.Fprintf(w, "%s", string(responseBytes))
	}
}

func approve_answer(w http.ResponseWriter, r *http.Request) {
	// get form fields
	answer_id := r.PostFormValue("answer_id")
	result, err := DB.Exec("update answers set approved=1 where id=?", answer_id)
	if err != nil {
		panic(err.Error())
	}
	rows, err := result.RowsAffected()
	if rows < 1 {
		panic(answer_id)
	}
	http.Redirect(w, r, "/control", http.StatusTemporaryRedirect)
}

func vote(w http.ResponseWriter, r *http.Request) {
	// get form fields
	answer_id := r.PostFormValue("answer_id")

	// get answer from DB
	queryResult := DB.QueryRow("select id, body, votes from answers where id=? and approved=1", answer_id)
	var a Answer
	queryResult.Scan(
		&a.Id,
		&a.Body,
		&a.Votes,
	)

	// Update vote count
	remainingVotes := get_remaining_votes(w, r)
	// fmt.Printf("remaining votes: %d", remainingVotes)

	if remainingVotes > 0 {
		result, err := DB.Exec("update answers set votes=? where id=?", a.Votes+1, answer_id)
		if err != nil {
			panic(err.Error())
		}
		rows, err := result.RowsAffected()
		if rows < 1 {
			panic(answer_id)
		}
		session, _ := store.Get(r, "session")
		session.Values["remaining_votes"] = remainingVotes - 1
		err = session.Save(r, w)
		if err != nil {
			fmt.Printf("error saving session")
		}
		// return response
		response := VoteResponse{
			Result:         "success",
			Votes:          a.Votes + 1,
			RemainingVotes: remainingVotes - 1,
		}
		responseJSON, err := json.Marshal(response)
		fmt.Fprint(w, string(responseJSON))
	} else {
		response := VoteResponse{
			Result: "no more votes",
			Votes:  a.Votes,
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprint(w, string(responseJSON))

	}

}

func post_answer(w http.ResponseWriter, r *http.Request) {

	// get form fields
	body := r.PostFormValue("body")
	author := r.PostFormValue("author")
	question_id := r.PostFormValue("question_id")

	// if the question has not been answered
	if !get_question_answered(w, r) {

		session, _ := store.Get(r, "session")

		session.Values["answered"] = true
		err := session.Save(r, w)
		if err != nil {
			fmt.Printf("error saving session")
		}

		// post answer

		_, err = DB.Exec("insert into answers(body, author, date_submitted, question_id) VALUES(?, ?, NOW(), ?)", body, author, question_id)
		if err != nil {
			panic(err.Error())
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func post_question(w http.ResponseWriter, r *http.Request) {

	// get form fields
	body := r.PostFormValue("body")

	// post question

	_, err := DB.Exec("insert into questions(body, date_submitted) VALUES(?, NOW())", body)
	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w, r, "/control", http.StatusTemporaryRedirect)
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
