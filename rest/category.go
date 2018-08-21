package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func createCategory(w http.ResponseWriter, r *http.Request, u *model.User) {

	category := &model.Category{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	err = getStorage().CreateCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategory(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	category := &model.Category{ID: id}

	err := getStorage().GetCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func updateCategory(w http.ResponseWriter, r *http.Request, u *model.User) {

	category := &model.Category{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&category)
	if err != nil {
		handleError(w, r, err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	category.ID = id

	err = getStorage().UpdateCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func deleteCategory(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	category := &model.Category{ID: id}

	err := getStorage().DeleteCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func getAllCategories(w http.ResponseWriter, r *http.Request, u *model.User) {

	categories, err := getStorage().GetCategories()
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
