package httpx

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

type decodeTarget struct {
	Name string `json:"name"`
}

func TestDecodeJSON_ValidBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"ivan"}`))

	var dst decodeTarget
	if err := DecodeJSON(rec, req, &dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Name != "ivan" {
		t.Fatalf("unexpected decoded value: %+v", dst)
	}
}

func TestDecodeJSON_UnknownFields_Rejected(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"ivan","extra":"field"}`))

	var dst decodeTarget
	if err := DecodeJSON(rec, req, &dst); err == nil {
		t.Fatalf("expected an error for unknown fields")
	}
}

func TestDecodeJSON_MultipleJSONObjects_Rejected(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"ivan"}{"name":"maria"}`))

	var dst decodeTarget
	if err := DecodeJSON(rec, req, &dst); err == nil {
		t.Fatalf("expected an error when the body contains more than one JSON value")
	}
}

func TestDecodeJSON_MalformedJSON_Rejected(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":`))

	var dst decodeTarget
	if err := DecodeJSON(rec, req, &dst); err == nil {
		t.Fatalf("expected an error for malformed JSON")
	}
}

func TestDecodeJSON_EmptyBody_Rejected(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(``))

	var dst decodeTarget
	if err := DecodeJSON(rec, req, &dst); err == nil {
		t.Fatalf("expected an error for an empty body")
	}
}

func TestWriteJSON_SetsStatusContentTypeAndBody(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(rec, 201, decodeTarget{Name: "ivan"})

	if rec.Code != 201 {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	var got decodeTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if got.Name != "ivan" {
		t.Fatalf("unexpected body: %+v", got)
	}
}
