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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	serveAddress        = ":8082"
	headerName   string = "Content-Type"
	headerValue  string = "application/json"
)

var errWrongRequest = errors.New("wrong request: check request headers")

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", signIn).Methods(http.MethodPost)
	r.HandleFunc("/users", getUsers).Methods(http.MethodGet)

	log.Println("server started")
	log.Fatal(http.ListenAndServe(serveAddress, r))
}

// Return all users.
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	res, err := storage.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		if _, err := w.Write([]byte(helpers.ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, res)
}

// Create new user.
func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerName, headerValue)

	user := models.User{
		ID:        "",
		FullName:  "",
		Email:     "",
		Password:  "",
		Active:    true,
		CreatedAt: "",
		UpdatedAt: "",
	}

	if checkRequestValidity(r) {
		w.WriteHeader(http.StatusNotAcceptable)

		err := fmt.Errorf("checkRequestValidity: %w", errWrongRequest)

		if _, err := w.Write([]byte(helpers.ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)

		if _, err := w.Write([]byte(helpers.ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	if err := user.ValidateUserField(); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)

		if _, err := w.Write([]byte(helpers.ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	if err := storage.CreateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		if _, err := w.Write([]byte(helpers.ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	if i, err := w.Write([]byte(helpers.ReturnOkJSON("Ok"))); err != nil {
		log.Println(i, err)
	}
}

// Check request validity.
func checkRequestValidity(r *http.Request) bool {
	acceptable := true

	for k, v := range r.Header {
		if k == headerName {
			for _, i := range v {
				if i == headerValue {
					acceptable = false
				}
			}
		}
	}

	return acceptable
}
