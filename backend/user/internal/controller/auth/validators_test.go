package auth

import (
	"strings"
	"testing"
)

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"minimum length (2)", "ab", true},
		{"maximum length (40)", strings.Repeat("a", 40), true},
		{"too short (1 char)", "a", false},
		{"too long (41 chars)", strings.Repeat("a", 41), false},
		{"empty", "", false},
		{"unicode counted by runes not bytes", "ая", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUsername(tt.username); got != tt.want {
				t.Errorf("isValidUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid simple email", "ivan@example.com", true},
		{"missing @", "ivan.example.com", false},
		{"empty", "", false},
		{"trailing garbage after valid address", "ivan@example.com, extra", false},
		{"display name form is rejected (address must equal input)", "Ivan <ivan@example.com>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"minimum length (8 bytes)", strings.Repeat("a", 8), true},
		{"maximum length (72 bytes)", strings.Repeat("a", 72), true},
		{"too short (7 bytes)", strings.Repeat("a", 7), false},
		{"too long (73 bytes)", strings.Repeat("a", 73), false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPassword(tt.password); got != tt.want {
				t.Errorf("isValidPassword(...) = %v, want %v", got, tt.want)
			}
		})
	}
}
