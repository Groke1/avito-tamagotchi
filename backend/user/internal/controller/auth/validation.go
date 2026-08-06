package auth

import "net/mail"

func isValidUsername(username string) bool {
	length := len([]rune(username))
	return length >= 2 && length <= 40
}

func isValidEmail(email string) bool {
	parsed, err := mail.ParseAddress(email)
	return err == nil && parsed.Address == email
}

func isValidPassword(password string) bool {
	length := len([]byte(password))
	return length >= 8 && length <= 72
}
