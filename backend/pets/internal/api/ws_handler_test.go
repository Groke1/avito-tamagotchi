package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/websocket"
)

func TestWSHandlerCreateWSTicket(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		handler := NewWSHandler(websocket.NewUserManager(), websocket.NewTicketManager(), nil)

		req := httptest.NewRequest(http.MethodPost, "/ws-ticket", nil)
		rec := httptest.NewRecorder()

		handler.CreateWSTicket(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("success", func(t *testing.T) {
		ticketManager := websocket.NewTicketManager()
		handler := NewWSHandler(websocket.NewUserManager(), ticketManager, nil)

		req := httptest.NewRequest(http.MethodPost, "/ws-ticket", nil)
		req = req.WithContext(context.WithValue(req.Context(), userIDKey, "user-1"))
		rec := httptest.NewRecorder()

		handler.CreateWSTicket(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var body struct {
			Ticket string `json:"ticket"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Ticket == "" {
			t.Fatal("ticket is empty")
		}

		userID, ok := ticketManager.ValidateAndDelete(body.Ticket)
		if !ok {
			t.Fatal("ticket was not registered in ticket manager")
		}
		if userID != "user-1" {
			t.Errorf("userID = %q, want %q", userID, "user-1")
		}
	})
}
