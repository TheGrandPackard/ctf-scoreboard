package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	// TODO Authorization - Need admin for create/update/delete. Enforce with HTTP middleware and JWT passed in with request

	r.HandleFunc("/api/category", createCategory).Methods("POST")               // Create a new category
	r.HandleFunc("/api/category/{id:[0-9]+}", getCategory).Methods("GET")       // Get existing category by ID
	r.HandleFunc("/api/category/{id:[0-9]+}", updateCategory).Methods("POST")   // Update existing category by ID
	r.HandleFunc("/api/category/{id:[0-9]+}", deleteCategory).Methods("DELETE") // Delete existing category by ID
	r.HandleFunc("/api/categories", getAllCategories).Methods("GET")            // Get all categories

	r.HandleFunc("/api/question", createQuestion).Methods("POST")               // Create a new question
	r.HandleFunc("/api/question/{id:[0-9]+}", getQuestion).Methods("GET")       // Get existing question by ID
	r.HandleFunc("/api/question/{id:[0-9]+}", updateQuestion).Methods("POST")   // Update existing question by ID
	r.HandleFunc("/api/question/{id:[0-9]+}", deleteQuestion).Methods("DELETE") // Delete existing question by ID
	r.HandleFunc("/api/questions", getAllQuestions).Methods("GET")              // Get all questions

	r.HandleFunc("/api/scores/question/{id:[0-9]+}", getQuestionsScores).Methods("GET") // Get all scores for a given question by ID
	r.HandleFunc("/api/scores/user/{id:[0-9]+}", getUserScores).Methods("GET")          // Get all scores for a given user by ID
	r.HandleFunc("/api/scores", getAllScores).Methods("GET")                            // Get all scores
	r.HandleFunc("/api/answer", submitAnswer).Methods("POST")                           // Submit an answer for a user

	r.HandleFunc("/api/team", createTeam).Methods("POST")               // Create a new team
	r.HandleFunc("/api/team/{id:[0-9]+}", getTeam).Methods("GET")       // Get existing team by ID
	r.HandleFunc("/api/team/{id:[0-9]+}", updateTeam).Methods("POST")   // Update existing team by ID
	r.HandleFunc("/api/team/{id:[0-9]+}", deleteTeam).Methods("DELETE") // Delete existing team by ID
	r.HandleFunc("/api/teams", getAllTeams).Methods("GET")              // Get all teams

	r.HandleFunc("/api/user", createUser).Methods("POST")                        // Create a new user
	r.HandleFunc("/api/user/{id:[0-9]+}", getUser).Methods("GET")                // Get existing user by ID
	r.HandleFunc("/api/user/{id:[0-9]+}", updateUser).Methods("POST")            // Update existing user by ID
	r.HandleFunc("/api/user/{id:[0-9]+}", deleteUser).Methods("DELETE")          // Delete existing user by ID
	r.HandleFunc("/api/users", getAllUsers).Methods("GET")                       // Get all users
	r.HandleFunc("/api/user/login", loginUser).Methods("POST")                   // Authenicate exister user for login, returns authorization token
	r.HandleFunc("/api/user/changepassword", changeUserPassword).Methods("POST") // Change password for logged in user

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

type jsonError struct {
	Error string `json:"error"`
}

func handleError(w http.ResponseWriter, r *http.Request, e error) {
	err := jsonError{Error: e.Error()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
