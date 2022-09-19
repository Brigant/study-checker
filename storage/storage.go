package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"study-checker/models"
)

const (
	dbhost string = "localhost"
	dbport string = "5432"
	dbname string = "study"
	dbuser string = "ps_user"
	dbpass string = "SimplePass"
	dbType string = "postgres" // not mysql.
)

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

// Return all users string json.
func GetAllUsers() (string, error) {
	dataBase, err := conenctDatabase()
	if err != nil {
		return "", err
	}
	defer dataBase.Close()

	rows, err := dataBase.Query(`SELECT id, name, email, active, created_at, updated_at FROM "user"`)
	if err != nil {
		return "", fmt.Errorf("dataBase.Query(): %w", err)
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		user := models.User{
			ID:        "",
			FullName:  "",
			Email:     "",
			Password:  "",
			Active:    true,
			CreatedAt: "",
			UpdatedAt: "",
		}

		err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return "", fmt.Errorf("rows.Scan: %w", err)
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return "", fmt.Errorf("%w,", err)
	}

	res, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	return string(res), nil
}

// Create user in database.
func CreateUser(user models.User) error {
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
