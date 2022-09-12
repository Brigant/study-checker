package models

import (
	"errors"
	"net/mail"
	"unicode"
)

type User struct {
	Id       string
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}

// validate user's field
func (u *User) ValidateUserField() error {
	var (
		errBadFullName   = errors.New("wrong full name")
		errWrongEmail    = errors.New("wrong email")
		errWrongPassword = errors.New("wrong password strength")
	)

	if len(u.FullName) < 2 || len(u.FullName) > 256 || !isASCII(u.FullName) {
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

// helper function for ASCII belonging
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
