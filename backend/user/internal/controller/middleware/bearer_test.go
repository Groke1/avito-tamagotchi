package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{"valid header", "Bearer abc.def.ghi", "abc.def.ghi", true},
		{"case-insensitive scheme", "bearer abc.def.ghi", "abc.def.ghi", true},
		{"empty header", "", "", false},
		{"missing token", "Bearer", "", false},
		{"wrong scheme", "Basic abc.def.ghi", "", false},
		{"too many parts", "Bearer abc def", "", false},
		{"empty token after scheme", "Bearer ", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, ok := bearerToken(tt.header)
			if ok != tt.wantOK {
				t.Fatalf("expected ok=%v, got %v", tt.wantOK, ok)
			}
			if token != tt.wantToken {
				t.Fatalf("expected token %q, got %q", tt.wantToken, token)
			}
		})
	}
}

type fakeValidator struct {
	validateFn func(ctx context.Context, token string) (string, error)
}

func (f *fakeValidator) ValidateAccessToken(ctx context.Context, token string) (string, error) {
	return f.validateFn(ctx, token)
}

func TestRequireAccessToken_ValidToken_SetsUserIDAndCallsNext(t *testing.T) {
	validator := &fakeValidator{validateFn: func(ctx context.Context, token string) (string, error) { return "user-1", nil }}
	nextCalled := false
	var seenUserID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		seenUserID, _ = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := RequireAccessToken(validator)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
	if seenUserID != "user-1" {
		t.Fatalf("expected userID user-1 in context, got %q", seenUserID)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRequireAccessToken_MissingHeader_Returns401(t *testing.T) {
	validator := &fakeValidator{validateFn: func(ctx context.Context, token string) (string, error) {
		t.Fatalf("validator should not be called without a bearer token")
		return "", nil
	}}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called")
	})

	handler := RequireAccessToken(validator)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAccessToken_ValidatorError_Returns401(t *testing.T) {
	validator := &fakeValidator{validateFn: func(ctx context.Context, token string) (string, error) {
		return "", context.DeadlineExceeded
	}}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called")
	})

	handler := RequireAccessToken(validator)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAccessToken_ValidatorReturnsEmptyUserID_Returns401(t *testing.T) {
	validator := &fakeValidator{validateFn: func(ctx context.Context, token string) (string, error) { return "", nil }}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called")
	})

	handler := RequireAccessToken(validator)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUserIDFromContext_NotSet(t *testing.T) {
	_, ok := UserIDFromContext(context.Background())
	if ok {
		t.Fatalf("expected ok=false for a context without a userID")
	}
}
