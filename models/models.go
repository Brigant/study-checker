package models

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"
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
)

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
