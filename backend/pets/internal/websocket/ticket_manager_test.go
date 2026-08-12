package websocket

import (
	"testing"
)

func TestTicketManagerCreateAndValidate(t *testing.T) {
	tm := NewTicketManager()

	ticket, err := tm.CreateTicket("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ticket == "" {
		t.Fatal("ticket is empty")
	}

	userID, ok := tm.ValidateAndDelete(ticket)
	if !ok {
		t.Fatal("expected ticket to be valid")
	}
	if userID != "user-1" {
		t.Errorf("userID = %q, want %q", userID, "user-1")
	}
}

func TestTicketManagerUnknownTicket(t *testing.T) {
	tm := NewTicketManager()

	_, ok := tm.ValidateAndDelete("nonexistent-ticket")
	if ok {
		t.Error("expected unknown ticket to be invalid")
	}
}

func TestTicketManagerDistinctTicketsPerCall(t *testing.T) {
	tm := NewTicketManager()

	ticket1, _ := tm.CreateTicket("user-1")
	ticket2, _ := tm.CreateTicket("user-1")

	if ticket1 == ticket2 {
		t.Error("expected distinct tickets on repeated calls")
	}
}

func TestTicketManagerMultipleUsers(t *testing.T) {
	tm := NewTicketManager()

	ticketA, _ := tm.CreateTicket("user-a")
	ticketB, _ := tm.CreateTicket("user-b")

	userA, okA := tm.ValidateAndDelete(ticketA)
	userB, okB := tm.ValidateAndDelete(ticketB)

	if !okA || userA != "user-a" {
		t.Errorf("userA = %q, ok = %v, want %q, true", userA, okA, "user-a")
	}
	if !okB || userB != "user-b" {
		t.Errorf("userB = %q, ok = %v, want %q, true", userB, okB, "user-b")
	}
}
