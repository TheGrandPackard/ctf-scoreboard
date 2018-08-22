package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
	"github.com/thegrandpackard/ctf-scoreboard/storage/mariadb"
)

var storage *mariadb.Storage

// TODO: Move this to the config file
var mySigningKey = []byte("s3cr3t")

func getStorage() *mariadb.Storage {
	return storage
}

// InitializeRoutes - initializeRoutes
func InitializeRoutes(s *mariadb.Storage) (err error) {

	storage = s

	r := mux.NewRouter()

	r.HandleFunc("/api/category", logHandler(requireAdmin(createCategory))).Methods("POST")               // Create a new category
	r.HandleFunc("/api/category/{id:[0-9]+}", logHandler(requireUser(getCategory))).Methods("GET")        // Get existing category by ID
	r.HandleFunc("/api/category/{id:[0-9]+}", logHandler(requireAdmin(updateCategory))).Methods("POST")   // Update existing category by ID
	r.HandleFunc("/api/category/{id:[0-9]+}", logHandler(requireAdmin(deleteCategory))).Methods("DELETE") // Delete existing category by ID
	r.HandleFunc("/api/categories", logHandler(requireUser(getAllCategories))).Methods("GET")             // Get all categories

	r.HandleFunc("/api/question", logHandler(requireAdmin(createQuestion))).Methods("POST")               // Create a new question
	r.HandleFunc("/api/question/{id:[0-9]+}", logHandler(requireUser(getQuestion))).Methods("GET")        // Get existing question by ID
	r.HandleFunc("/api/question/{id:[0-9]+}", logHandler(requireAdmin(updateQuestion))).Methods("POST")   // Update existing question by ID
	r.HandleFunc("/api/question/{id:[0-9]+}", logHandler(requireAdmin(deleteQuestion))).Methods("DELETE") // Delete existing question by ID
	r.HandleFunc("/api/questions", logHandler(requireUser(getAllQuestions))).Methods("GET")               // Get all questions

	r.HandleFunc("/api/scores/question/{id:[0-9]+}", logHandler(requireUser(getQuestionsScores))).Methods("GET") // Get all scores for a given question by ID
	r.HandleFunc("/api/scores/user/{id:[0-9]+}", logHandler(requireUser(getUserScores))).Methods("GET")          // Get all scores for a given user by ID
	r.HandleFunc("/api/scores", logHandler(requireUser(getAllScores))).Methods("GET")                            // Get all scores
	r.HandleFunc("/api/answer", logHandler(requireUser(submitAnswer))).Methods("POST")                           // Submit an answer for a user

	r.HandleFunc("/api/team", logHandler(requireAdmin(createTeam))).Methods("POST")               // Create a new team
	r.HandleFunc("/api/team/{id:[0-9]+}", logHandler(requireUser(getTeam))).Methods("GET")        // Get existing team by ID
	r.HandleFunc("/api/team/{id:[0-9]+}", logHandler(requireAdmin(updateTeam))).Methods("POST")   // Update existing team by ID
	r.HandleFunc("/api/team/{id:[0-9]+}", logHandler(requireAdmin(deleteTeam))).Methods("DELETE") // Delete existing team by ID
	r.HandleFunc("/api/teams", logHandler(requireUser(getAllTeams))).Methods("GET")               // Get all teams

	r.HandleFunc("/api/user", logHandler(requireAdmin(createUser))).Methods("POST")                       // Create a new user
	r.HandleFunc("/api/user/{id:[0-9]+}", logHandler(requireUser(getUser))).Methods("GET")                // Get existing user by ID
	r.HandleFunc("/api/user/{id:[0-9]+}", logHandler(requireAdmin(updateUser))).Methods("POST")           // Update existing user by ID
	r.HandleFunc("/api/user/{id:[0-9]+}", logHandler(requireAdmin(deleteUser))).Methods("DELETE")         // Delete existing user by ID
	r.HandleFunc("/api/users", logHandler(requireUser(getAllUsers))).Methods("GET")                       // Get all users
	r.HandleFunc("/api/user/login", logHandler(loginUser)).Methods("POST")                                // Authenicate exister user for login, returns authorization token
	r.HandleFunc("/api/user/changepassword", logHandler(requireUser(changeUserPassword))).Methods("POST") // Change password for logged in user

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	return
}

type authenticatedHandlerFunc (func(w http.ResponseWriter, r *http.Request, u *model.User))

// HTTP Logging and Authentication
func logHandler(next func(w http.ResponseWriter, r *http.Request, u *model.User)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := &model.User{Username: "-"}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			next(w, r, nil)
			logRequest(r, user)
			return
		}

		// Parse token from header
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return mySigningKey, nil
		})
		if err != nil {
			handleError(w, r, err)
			logRequest(r, user)
			return
		}

		// Get user from token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//log.Printf("ID: %+v, Username: %+v, expiration: %+v", claims["id"], claims["username"], claims["expiration"])

			userID, ok := claims["id"].(float64) // This is a float64 even though it is stored as a uint64 when it is created
			if ok {
				user = &model.User{ID: uint64(userID)}
				err = getStorage().GetUser(user)
				if err != nil {
					handleError(w, r, err)
					logRequest(r, user)
					return
				}
				next(w, r, user)
				logRequest(r, user)
				return
			} else {
				handleError(w, r, errors.New("Invalid token"))
				logRequest(r, user)
				return
			}
		} else {
			handleError(w, r, err)
			logRequest(r, user)
			return
		}
	})
}

func logRequest(r *http.Request, user *model.User) {
	// TODO: Implement new response writer in order to get response code and response size: https://www.reddit.com/r/golang/comments/7p35s4/how_do_i_get_the_response_status_for_my_middleware/
	//127.0.0.1 frank "GET /apache_pb.gif HTTP/1.0" 200 2326 "Mozilla/4.08 [en] (Win98; I ;Nav)"
	log.Printf("%s %s \"%s\" %d %d \"%s\"",
		r.RemoteAddr,
		user.Username,
		r.URL,
		0, //"response_code",
		0, //"response_size",
		r.UserAgent())
}

// Authorization
func requireUser(next authenticatedHandlerFunc) authenticatedHandlerFunc {
	return authenticatedHandlerFunc(func(w http.ResponseWriter, r *http.Request, u *model.User) {
		if u != nil {
			next(w, r, u)
		} else {
			handleError(w, r, errors.New("User Authorization Required"))
			return
		}
	})
}

func requireAdmin(next authenticatedHandlerFunc) authenticatedHandlerFunc {
	return authenticatedHandlerFunc(func(w http.ResponseWriter, r *http.Request, u *model.User) {
		if u != nil && u.Admin {
			next(w, r, u)
		} else {
			handleError(w, r, errors.New("Admin Authorization Required"))
			return
		}
	})
}

// Basic error handling
type jsonError struct {
	Error string `json:"error"`
}

func handleError(w http.ResponseWriter, r *http.Request, e error) {
	err := jsonError{Error: e.Error()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
