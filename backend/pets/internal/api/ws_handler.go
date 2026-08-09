package api

import (
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/websocket"
)

type WSHandler struct {
	client        *websocket.ClientManager
	ticketManager *websocket.TicketManager
	service       *service.PetService
}

func NewWSHandler(client *websocket.ClientManager, ticketManager *websocket.TicketManager, petService *service.PetService) *WSHandler {
	return &WSHandler{client: client, ticketManager: ticketManager, service: petService}
}

func (wsh *WSHandler) CreateWSTicket(w http.ResponseWriter, r *http.Request) {
	userID, err := UserIDFromContext(r.Context())
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	ticket, err := wsh.ticketManager.CreateTicket(userID)
	if err != nil {
		writeError(w, ErrInternalError)
		return
	}

	writeJSONResponse(w, http.StatusOK, struct {
		Ticket string `json:"ticket"`
	}{
		Ticket: ticket,
	})
}

func (wsh *WSHandler) OpenWSConn(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		writeError(w, ErrUnauthorized)
		return
	}

	userID, ok := wsh.ticketManager.ValidateAndDelete(ticket)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	wsh.client.ConnectClient(w, r, userID, wsh.service)
}
