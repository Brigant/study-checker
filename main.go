package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	serveAddress        = ":8082"
	headerName   string = "Content-Type"
	headerValue  string = "application/json"
	dbhost       string = "localhost"
	dbport       string = "5432"
	dbname       string = "study"
	dbuser       string = "ps_user"
	dbpass       string = "SimplePass"
	dbType       string = "postgres" // not mysql.
)

type User struct {
	ID        string
	FullName  string `json:"fullname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var (
	errBadFullName   = errors.New("wrong full name")
	errWrongEmail    = errors.New("wrong email")
	errWrongPassword = errors.New("wrong password strength")
	errWrongRequest  = errors.New("wrong request")
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", signIn).Methods(http.MethodPost)

	log.Println("server started")
	log.Fatal(http.ListenAndServe(serveAddress, r))
}

// Connect to database.
func conenctDatabase() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)

	dataBase, err := sql.Open(dbType, psqlconn)
	if err != nil {
		return nil, fmt.Errorf("conectDatabase(): %w", err)
	}

	return dataBase, nil
}

// For errors in User structer.
func userError(e error, msg string) error {
	return fmt.Errorf("%w: %s", e, msg)
}

// Validate user's field.
func (u *User) ValidateUserField() error {
	fullName := strings.TrimSpace("\t " + u.FullName + "\n ")

	u.FullName = fullName

	u.Email = strings.TrimSpace("\t " + u.Email + "\n ")

	if len(u.FullName) < 2 || len(u.FullName) > 256 {
		return userError(errBadFullName, "Full name can contains more then 2 characters and less then 256")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil || len(u.Email) > 256 || !isASCII(u.Email) {
		return userError(errWrongEmail, "Email should be  more 2 and less 256 and in ASCII characters")
	}

	if len(u.Password) < 9 || len(u.Password) > 256 {
		return userError(errWrongPassword, "Password can contains more then 2 characters and less then 256")
	}

	return nil
}

// Helper function for ASCII belonging.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// Create user in database.
func CreateUser(user User) error {
	dataBase, err := conenctDatabase()
	if err != nil {
		return fmt.Errorf("%w,", err)
	}
	defer dataBase.Close()

	dbRequest := `INSERT INTO public.user (name, email, active, password) VALUES($1, $2, $3, md5($4));`

	_, err = dataBase.Exec(dbRequest, user.FullName, user.Email, user.Active, user.Password)
	if err != nil {
		return fmt.Errorf("%w,", err)
	}

	return nil
}

// Create new user.
func signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerName, headerValue)

	user := User{
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

		if _, err := w.Write([]byte(ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)

		if _, err := w.Write([]byte(ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	if err := user.ValidateUserField(); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)

		if _, err := w.Write([]byte(ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	if err := CreateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		if _, err := w.Write([]byte(ReturnErrorJSON(err))); err != nil {
			log.Println(err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	if i, err := w.Write([]byte(ReturnOkJSON("Ok"))); err != nil {
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

// Just convert error to json.
func ReturnErrorJSON(err error) string {
	type JError struct {
		Result     string
		Time       string
		ErrDetails string
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	jerror := JError{
		Result:     "Error",
		Time:       currentTime,
		ErrDetails: err.Error(),
	}

	res, e := json.Marshal(jerror)
	if e != nil {
		log.Println(e)
	}

	s := string(res)

	return s
}

// Just return pretty "OK" in json.
func ReturnOkJSON(result string) string {
	type JError struct {
		Result string
		Time   string
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	je := JError{
		Result: result,
		Time:   currentTime,
	}

	res, e := json.Marshal(je)
	if e != nil {
		log.Println(e)
	}

	s := string(res)

	return s
}
