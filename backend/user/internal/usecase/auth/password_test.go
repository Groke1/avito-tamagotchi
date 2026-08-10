package auth

import "testing"

func TestHashPassword_And_CheckPasswordHash_Roundtrip(t *testing.T) {
	encoded, err := hashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !checkPasswordHash("correct-horse-battery-staple", encoded) {
		t.Fatalf("expected correct password to match its own hash")
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	encoded, err := hashPassword("correct-password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if checkPasswordHash("wrong-password", encoded) {
		t.Fatalf("expected wrong password not to match")
	}
}

func TestHashPassword_DifferentSaltsProduceDifferentHashes(t *testing.T) {
	h1, err := hashPassword("same-password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := hashPassword("same-password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h1 == h2 {
		t.Fatalf("expected different random salts to produce different encoded hashes")
	}
	if !checkPasswordHash("same-password", h1) || !checkPasswordHash("same-password", h2) {
		t.Fatalf("both hashes should validate against the original password")
	}
}

func TestCheckPasswordHash_MalformedInputs(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{"empty string", ""},
		{"wrong number of segments", "pbkdf2-sha256$100000$salt"},
		{"wrong algorithm tag", "bcrypt$100000$c2FsdA$aGFzaA"},
		{"non-numeric iterations", "pbkdf2-sha256$abc$c2FsdA$aGFzaA"},
		{"zero iterations", "pbkdf2-sha256$0$c2FsdA$aGFzaA"},
		{"invalid base64 salt", "pbkdf2-sha256$100000$not-valid-base64!!$aGFzaA"},
		{"invalid base64 hash", "pbkdf2-sha256$100000$c2FsdA$not-valid-base64!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if checkPasswordHash("any-password", tt.encoded) {
				t.Fatalf("expected malformed encoded hash to fail validation")
			}
		})
	}
}
