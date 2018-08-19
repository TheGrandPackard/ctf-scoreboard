package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// TODO Authorization - Needs admin for create/update/delete

func createCategory(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	category := &model.Category{Name: r.Form.Get("name")}
	err := getStorage().CreateCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseUint(idString, 10, 64)
	category := &model.Category{ID: id}
	err = getStorage().GetCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)

}

func updateCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseUint(idString, 10, 64)
	r.ParseForm()
	category := &model.Category{ID: id, Name: r.PostForm.Get("name")}
	fmt.Printf("%+v", category)
	err = getStorage().UpdateCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseUint(idString, 10, 64)
	category := &model.Category{ID: id}
	err = getStorage().DeleteCategory(category)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {

	categories, err := getStorage().GetCategories()
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
