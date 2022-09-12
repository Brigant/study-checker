package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"study-checker/models"
	"study-checker/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const serveAddress = ":8082"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", signIn).Methods(http.MethodPost)
	r.HandleFunc("/users", getUsers).Methods(http.MethodGet)

	log.Println("server started")
	log.Fatal(http.ListenAndServe(serveAddress, r))
}

// return all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	res, err := storage.GetAllUsers()
	if err != nil {
		http.Error(w, "Internak server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(res))
}

// create new user
func signIn(w http.ResponseWriter, r *http.Request) {
	const (
		headerName  string = "Content-Type"
		headerValue string = "application/json"
	)

	var (
		user        models.User
		acceptable  bool
		isEmailInDB bool
	)

	//check request validity
	for k, v := range r.Header {
		if k == headerName {
			for _, i := range v {
				if i == headerValue {
					acceptable = true
				}
			}
		}
	}
	if !acceptable {
		http.Error(w, "Something wrong in your request", http.StatusNotAcceptable)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Id = uuid.New().String()
	user.Active = true
	if err := user.ValidateUserField(); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	emails, err := storage.GetEmails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, emailInDB := range emails {
		if emailInDB == user.Email {
			isEmailInDB = true
		}
	}

	if isEmailInDB {
		fmt.Fprintln(w, "User already exists")
		return
	}

	if err := storage.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "user: %+v", user)
}
