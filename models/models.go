package models

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
)

type User struct {
	Id         string
	FullName   string `json:"fullname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Active     bool   `json:"active"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

// validate user's field
func (u *User) ValidateUserField() error {
	var (
		errBadFullName   = errors.New("wrong full name")
		errWrongEmail    = errors.New("wrong email")
		errWrongPassword = errors.New("wrong password strength")
	)
	u.FullName = strings.TrimSpace("\t " + u.FullName + "\n ")
	u.Email = strings.TrimSpace("\t " + u.Email + "\n ")
	if len(u.FullName) < 2 || len(u.FullName) > 256 {
		return errBadFullName
	}
	if _, err := mail.ParseAddress(u.Email); err != nil || len(u.Email) > 256 || !isASCII(u.Email) {
		return errWrongEmail
	}
	if len(u.Password) < 9 || len(u.Password) > 256 {
		return errWrongPassword
	}
	return nil
}

// retur user JSON bytes
// func (u User) ReturneInJSON() string {
// 	u.Password = "******"
// 	str, err := json.Marshal(u)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
// 	return string(str)
// }

// helper function for ASCII belonging
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
