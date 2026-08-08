package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type TicketManager struct {
	tickets sync.Map
}

func NewTicketManager() *TicketManager {
	return &TicketManager{}
}

func (tm *TicketManager) CreateTicket(userID string) (string, error) {
	// TODO сделать более защищенную генерацию
	bytes := make([]byte, 16)
	rand.Read(bytes)

	ticket := hex.EncodeToString(bytes)

	tm.tickets.Store(ticket, userID)

	time.AfterFunc(2*time.Minute, func() {
		tm.tickets.Delete(userID)
	})

	return ticket, nil
}

func (tm *TicketManager) ValidateAndDelete(ticket string) (string, bool) {
	userID, ok := tm.tickets.Load(ticket)
	if !ok {
		return "", false
	}

	return userID.(string), true
}
