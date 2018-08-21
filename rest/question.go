package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func createQuestion(w http.ResponseWriter, r *http.Request, u *model.User) {

	question := &model.Question{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	err = getStorage().CreateQuestion(question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func getQuestion(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	question := &model.Question{ID: id}

	err := getStorage().GetQuestion(question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func updateQuestion(w http.ResponseWriter, r *http.Request, u *model.User) {

	question := &model.Question{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		handleError(w, r, err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	question.ID = id

	err = getStorage().UpdateQuestion(question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func deleteQuestion(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	question := &model.Question{ID: id}

	err := getStorage().DeleteQuestion(question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func getAllQuestions(w http.ResponseWriter, r *http.Request, u *model.User) {

	questions, err := getStorage().GetQuestions()
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}
