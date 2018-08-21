package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/thegrandpackard/ctf-scoreboard/model"
	"github.com/thegrandpackard/ctf-scoreboard/storage/mariadb"
)

var storage *mariadb.Storage

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
		log.Printf("Handled HTTP Request: %+v", r)
		// TODO: Make like apache logging
		// TODO: Authenticate user to add to logging and also pass through for authorization to endpoints
		u := &model.User{}
		next(w, r, u)
	})
}

// Authorization
func requireUser(next authenticatedHandlerFunc) authenticatedHandlerFunc {
	return authenticatedHandlerFunc(func(w http.ResponseWriter, r *http.Request, u *model.User) {
		// TODO: Test for valid user
		next(w, r, u)
	})
}

func requireAdmin(next authenticatedHandlerFunc) authenticatedHandlerFunc {
	return authenticatedHandlerFunc(func(w http.ResponseWriter, r *http.Request, u *model.User) {
		// TODO: Test for user is admin
		next(w, r, u)
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
