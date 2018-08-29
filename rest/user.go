package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/config"
	"github.com/thegrandpackard/ctf-scoreboard/model"
	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request, u *model.User) {

	user := &model.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	// bcrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		handleError(w, r, err)
		return
	}
	user.Password = string(hash)

	// store the user
	err = getStorage().CreateUser(user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	user := &model.User{ID: id}

	err := getStorage().GetUser(user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request, u *model.User) {

	user := &model.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		handleError(w, r, err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	user.ID = id

	err = getStorage().UpdateUser(user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request, u *model.User) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	user := &model.User{ID: id}

	err := getStorage().DeleteUser(user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func getAllUsers(w http.ResponseWriter, r *http.Request, u *model.User) {

	categories, err := getStorage().GetCategories()
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func changeUserPassword(w http.ResponseWriter, r *http.Request, u *model.User) {

	user := &model.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		handleError(w, r, err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	user.ID = id

	if u.ID == user.ID {

		// bcrypt the password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
		if err != nil {
			handleError(w, r, err)
			return
		}
		user.Password = string(hash)

		err = getStorage().UpdateUserPassword(user)
		if err != nil {
			handleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		handleError(w, r, errors.New("Incorrect answer"))
		return
	}

}

func loginUser(w http.ResponseWriter, r *http.Request, u *model.User) {

	user := &model.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	dbUser := &model.User{Username: user.Username}
	err = getStorage().GetUserAuthentication(dbUser)
	if err != nil && err == sql.ErrNoRows {
		handleError(w, r, errors.New("Invalid username or password"))
		return
	} else if err != nil {
		handleError(w, r, err)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) == nil {
		/* Create the token */
		token := jwt.New(jwt.SigningMethodHS256)

		/* Create a map to store our claims */
		claims := token.Claims.(jwt.MapClaims)

		/* Set token claims */
		claims["id"] = dbUser.ID
		claims["username"] = dbUser.Username
		claims["expiration"] = time.Now().Add(time.Hour * 24).Unix()

		/* Sign the token with our secret */
		configuration := config.LoadConfig()
		tokenString, err := token.SignedString([]byte(configuration.AuthSecret))
		if err != nil {
			handleError(w, r, err)
			return
		}

		/* Finally, write the token to the browser window */
		w.Write([]byte("{ \"token\": \"" + tokenString + "\" }"))
	} else {
		handleError(w, r, errors.New("Invalid username or password"))
		return
	}
}
