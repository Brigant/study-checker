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
	serveAddress           = ":8082"
	headerName      string = "Content-Type"
	headerValueJSON string = "application/json"
	dbhost          string = "localhost"
	dbport          string = "5432"
	dbname          string = "study"
	dbuser          string = "ps_user"
	dbpass          string = "SimplePass"
	dbType          string = "postgres" // not mysql.
)

var (
	errRequestValidation   = errors.New("checkRequestValidity")
	errUserNameValidation  = errors.New("full name should contain more then 2 characters and less then 256")
	errUserEmailValidation = errors.New("email should be  more 2 and less 256 and in ASCII characters")
	errUserPassValidation  = errors.New("password can contains more then 2 characters and less then 256")
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

type APIServer struct {
	address string
	store   *Store
	router  *mux.Router
}

type Store struct {
	dbType     string
	connString string
	db         *sql.DB
}

// Init the server.
func NewAPIServer() *APIServer {
	return &APIServer{
		address: serveAddress,
		router:  mux.NewRouter(),
	}
}

// Init the storage with database.
func NewStore(conString string, dbType string) *Store {
	return &Store{
		dbType:     dbType,
		connString: conString,
	}
}

// Starts the API server.
func (s *APIServer) Start() error {
	s.createRoutes()

	if err := s.configureStore(); err != nil {
		return err
	}

	log.Println("SerSver started")

	if err := http.ListenAndServe(s.address, s.router); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (s *APIServer) configureStore() error {
	psqlconn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)

	store := NewStore(psqlconn, dbType)

	if err := store.Open(); err != nil {
		return err
	}

	s.store = store

	return nil
}

func (s *APIServer) createRoutes() {
	s.router.HandleFunc("/users", s.handleSignIn).Methods(http.MethodPost)
}

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerName, headerValueJSON)

	if checkRequestValidity(r, headerName, headerValueJSON) {
		s.httpErrorJSON(w, http.StatusNotAcceptable, errRequestValidation)

		return
	}

	user := User{}
	user.Active = true

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.httpErrorJSON(w, http.StatusUnprocessableEntity, err)

		return
	}

	if err := user.validateUserField(); err != nil {
		s.httpErrorJSON(w, http.StatusUnprocessableEntity, err)

		return
	}

	if err := s.store.CreateUser(&user); err != nil {
		s.httpErrorJSON(w, http.StatusBadRequest, err)

		return
	}

	s.returnHTTPOkJSON(w, user)
}

// Helper function for http respond with error in json format.
func (s *APIServer) httpErrorJSON(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	some := struct {
		Result     string
		Time       string
		ErrDetails string
	}{
		Result:     "Error",
		Time:       currentTime,
		ErrDetails: err.Error(),
	}

	res, e := json.Marshal(some)
	if e != nil {
		log.Println(e)
	}

	if _, err := w.Write([]byte(string(res))); err != nil {
		log.Println(err)
	}
}

// Helper function for http respond with some object in json format.
func (s *APIServer) returnHTTPOkJSON(w http.ResponseWriter, u User) {
	w.WriteHeader(http.StatusOK)

	res, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}

	if _, err := w.Write([]byte(string(res))); err != nil {
		log.Println(err)
	}
}

// Creats user in database.
func (store *Store) CreateUser(user *User) error {
	dbRequest := `INSERT INTO public.user (name, email, active, password) VALUES($1, $2, $3, md5($4))
		RETURNING id, created_at, updated_at;`

	if err := store.db.QueryRow(
		dbRequest,
		&user.FullName,
		&user.Email,
		&user.Active,
		&user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

// Open connection to database.
func (store *Store) Open() error {
	db, err := sql.Open(store.dbType, store.connString)
	if err != nil {
		return fmt.Errorf("sql.Open(): %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db.Ping(): %w", err)
	}

	store.db = db

	log.Println("Connected to database")

	return nil
}

// Close connection to database.
func (store *Store) Close() {
	store.db.Close()
}

func (u *User) validateUserField() error {
	fullName := strings.TrimSpace("\t " + u.FullName + "\n ")

	u.FullName = fullName

	u.Email = strings.TrimSpace("\t " + u.Email + "\n ")

	if len(u.FullName) < 2 || len(u.FullName) > 256 {
		return errUserNameValidation
	}

	if _, err := mail.ParseAddress(u.Email); err != nil || len(u.Email) > 256 || !isASCII(u.Email) {
		return errUserEmailValidation
	}

	if len(u.Password) < 9 || len(u.Password) > 256 {
		return errUserPassValidation
	}

	return nil
}

func main() {
	server := NewAPIServer()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// helper func which check header on existing of the value.
func checkRequestValidity(r *http.Request, headerName, headerValue string) bool {
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

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}
