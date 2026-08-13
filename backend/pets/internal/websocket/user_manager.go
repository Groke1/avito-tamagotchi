package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
	"github.com/gorilla/websocket"
)

const (
	BroadcastLeaderboardEventCooldown = 1 * time.Minute
	PetUpdateEventCooldown            = 30 * time.Second
)

var wsUpgrader = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool {
	return true
}}

type UserManager struct {
	users                map[string]*UserConnections
	mu                   sync.RWMutex
	updateLeaderboardSig chan struct{}
}

func NewUserManager() *UserManager {
	um := &UserManager{users: make(map[string]*UserConnections), updateLeaderboardSig: make(chan struct{}, 1)}

	go um.leaderboadDebounce()

	return um
}

func (um *UserManager) ConnectClient(w http.ResponseWriter, r *http.Request, userID string, petService *service.PetService) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS ERROR] Upgrade failed for userID '%s': %v", userID, err)
		return
	}

	client := &ClientConnection{
		userID: userID,
		conn:   conn,
	}

	um.mu.Lock()

	if um.users[userID] == nil {
		userCtx, userCancel := context.WithCancel(context.Background())

		um.users[userID] = &UserConnections{
			clients: make(map[*ClientConnection]struct{}),
			cancel:  userCancel,
		}

		go um.startUpdates(userCtx, userID, petService)
	}

	um.users[userID].clients[client] = struct{}{}

	activeCount := len(um.users)
	um.mu.Unlock()

	log.Printf("[WS] Client connected successfully: userID='%s'. Total active users: %d", userID, activeCount)

	go um.listenAndServe(userID, client)
}

func (um *UserManager) startUpdates(ctx context.Context, userID string, petService *service.PetService) {
	ticker := time.NewTicker(PetUpdateEventCooldown)
	defer ticker.Stop()

	log.Printf("[WS TICKER] Started update loop for userID: '%s'", userID)

	for {
		select {
		case <-ticker.C:
			log.Printf("[WS TICKER] Tick fired for userID: '%s'", userID)

			pet, changed, err := petService.RecalculateState(ctx, userID)
			if err != nil {
				log.Printf("[WS TICKER ERROR] RecalculateState failed for userID '%s': %v", userID, err)
				continue
			}

			log.Printf("[WS TICKER] Pet retrieved. Changed=%t for userID: '%s'", changed, userID)
			if !changed {
				continue
			}

			petResponse := PetResponse{
				ID:          pet.ID,
				Name:        pet.Name,
				Level:       pet.Level,
				XP:          pet.XP,
				NextLevelXP: pet.NextLevelXP,
				Satiety:     pet.Satiety,
				Happiness:   pet.Happiness,
			}

			um.SendToClient(userID, domain.EventPetUpdated, petResponse)

		case <-ctx.Done():
			log.Printf("[WS TICKER] Context done, stopping loop for userID: '%s'", userID)
			return
		}
	}
}

func (um *UserManager) listenAndServe(userID string, client *ClientConnection) {
	defer func() {
		um.mu.Lock()

		user := um.users[userID]

		delete(user.clients, client)

		if len(user.clients) == 0 {
			user.cancel()
			delete(um.users, userID)
		}
		um.mu.Unlock()

		client.Close()
		log.Printf("[WS] User '%s' disconnected with some conn", userID)
	}()

	for {
		if _, _, err := client.ReadMessage(); err != nil {
			return
		}
	}
}

func (um *UserManager) BroadcastLeaderboard() {
	log.Println("[WS BROADCAST] BroadcastLeaderboardUpdate called from service")
	select {
	case um.updateLeaderboardSig <- struct{}{}:
		log.Println("[WS BROADCAST] Signal pushed to updateLeaderboardSig channel")
	default:
		log.Println("[WS BROADCAST] Channel full, signal skipped")
	}
}

func (um *UserManager) leaderboadDebounce() {
	log.Println("[WS DEBOUNCER] Debouncer started")

	ticker := time.NewTicker(BroadcastLeaderboardEventCooldown)
	defer ticker.Stop()

	for range um.updateLeaderboardSig {
		log.Println("[WS DEBOUNCER] Received signal from channel, waiting ticker")
		<-ticker.C

		event := Event{
			EventType: "leaderboard.position_updated",
		}

		um.mu.RLock()
		clients := make([]*ClientConnection, 0, len(um.users))
		for _, user := range um.users {
			for client := range user.clients {
				clients = append(clients, client)
			}
		}
		um.mu.RUnlock()

		for _, client := range clients {
			err := client.WriteJSON(event)
			if err != nil {
				log.Printf("[WS DEBOUNCER ERROR] Failed to send to userID '%s': %v", client.userID, err)
			}
		}
	}
}

func (um *UserManager) SendToClient(userID string, eventType domain.WsEvent, v any) bool {
	log.Println("[WS CLIENT] sender started")

	event := Event{
		EventType: eventType,
		Payload:   v,
	}

	um.mu.RLock()

	user, exists := um.users[userID]
	if !exists {
		um.mu.RUnlock()
		log.Printf("[WS CLIENT] User '%s' is offline, event skipped", userID)
		return false
	}

	clients := make([]*ClientConnection, 0, len(user.clients))
	for client := range user.clients {
		clients = append(clients, client)
	}

	um.mu.RUnlock()

	sent := false
	for _, client := range clients {
		if err := client.WriteJSON(event); err != nil {
			log.Printf("[WS CLIENT ERROR] Failed to send to userID '%s': %v", userID, err)
			continue
		}
		sent = true
	}

	return sent
}
