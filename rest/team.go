package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func createTeam(w http.ResponseWriter, r *http.Request, u *model.User) {

	team := &model.Team{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&team)
	if err != nil {
		handleError(w, r, err)
		return
	}

	err = getStorage().CreateTeam(team)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func getTeam(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	team := &model.Team{ID: id}

	err := getStorage().GetTeam(team)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func updateTeam(w http.ResponseWriter, r *http.Request, u *model.User) {

	team := &model.Team{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&team)
	if err != nil {
		handleError(w, r, err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	team.ID = id

	err = getStorage().UpdateTeam(team)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func deleteTeam(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	team := &model.Team{ID: id}

	err := getStorage().DeleteTeam(team)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func getAllTeams(w http.ResponseWriter, r *http.Request, u *model.User) {

	categories, err := getStorage().GetCategories()
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
