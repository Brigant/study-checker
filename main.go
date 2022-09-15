package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"study-checker/helpers"
	"study-checker/models"
	"study-checker/storage"
	"time"

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
	w.Header().Set("content-type", "application/json")
	res, err := storage.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(helpers.ReturnErrorJson(err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(res))
}

// create new user
func signIn(w http.ResponseWriter, r *http.Request) {
	const (
		headerName  string = "Content-Type"
		headerValue string = "application/json"
	)
	w.Header().Set(headerName, headerValue)

	var (
		user            models.User
		acceptable      bool
		isEmailInDB     bool
		errWrongRequest = errors.New("wrong request: check headers or request body")
		errUserExists   = errors.New("user already exists")
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
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(helpers.ReturnErrorJson(errWrongRequest)))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(helpers.ReturnErrorJson(err)))
		return
	}

	// create and fill struct User
	user.Id = uuid.New().String()
	user.Active = true
	user.Created = time.Now().Format("2006-01-02 15:04:05")
	user.Modified = time.Now().Format("2006-01-02 15:04:05")
	if err := user.ValidateUserField(); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(helpers.ReturnErrorJson(err)))
		return
	}

	emails, err := storage.GetEmails()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(helpers.ReturnErrorJson(err)))
		return
	}
	for _, emailInDB := range emails {
		if emailInDB == user.Email {
			isEmailInDB = true
		}
	}

	if isEmailInDB {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(helpers.ReturnErrorJson(errUserExists)))
		return
	}

	if err := storage.CreateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(helpers.ReturnErrorJson(errUserExists)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(user.ReturnJsonString()))
}
