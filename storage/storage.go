package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"study-checker/models"

	_ "github.com/lib/pq"
)

const (
	dbhost string = "127.0.0.1"
	dbport string = "49153"
	dbname string = "postgres"
	dbuser string = "checker"
	dbpass string = "password"
)

// return all users string json
func GetAllUsers() (string, error) {
	psqlconn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, email, active, password FROM "user"`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var users = []models.User{}
	for rows.Next() {
		var user = models.User{}

		err := rows.Scan(&user.Id, &user.FullName, &user.Email, &user.Active, &user.Password)
		if err != nil {
			return "", err
		}

		users = append(users, user)
	}
	res, err := json.Marshal(users)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// return all user's emails
func GetEmails() ([]string, error) {
	psqlconn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT email FROM "user"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails = []string{}
	for rows.Next() {
		var email string

		err := rows.Scan(&email)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// create user in database
func CreateUser(u models.User) error {
	psqlconn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbRequest := `INSERT INTO public."user" (id, "name", email, active, "password") VALUES($1, $2, $3, $4, $5);`
	_, err = db.Exec(dbRequest, u.Id, u.FullName, u.Email, u.Active, u.Password)
	if err != nil {
		return err
	}

	return nil
}
