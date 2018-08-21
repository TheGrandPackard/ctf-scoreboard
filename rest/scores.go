package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func getQuestionsScores(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	question := &model.Question{ID: id}

	scores, err := getStorage().GetQuestionScores(question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

func getUserScores(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	user := &model.User{ID: id}

	scores, err := getStorage().GetUserScores(user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

func getAllScores(w http.ResponseWriter, r *http.Request, u *model.User) {

	scores, err := getStorage().GetScores()
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

func submitAnswer(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Parse answer submitted by user
	answer := &model.Question{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&answer)
	if err != nil {
		handleError(w, r, err)
		return
	}

	// Get question from datbase to verify
	question := &model.Question{ID: answer.ID}
	err = getStorage().GetQuestion(question)
	if err != nil {
		handleError(w, r, err)
		return
	}

	// Check if answer is correct
	if question.Answer == answer.Answer {
		score := &model.Score{}
		score.User.ID = 0 // TODO Get from JWT
		score.Question.ID = question.ID
		err = getStorage().CreateScore(score)
		if err != nil {
			handleError(w, r, errors.New("Incorrect answer"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(score)
	} else {
		handleError(w, r, errors.New("Incorrect answer"))
		return
	}
}
