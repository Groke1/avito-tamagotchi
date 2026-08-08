package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

type ClientManager struct {
	clients              map[string]*ClientConnection
	mu                   sync.RWMutex
	updateLeaderboardSig chan struct{}
}

func NewClient() *ClientManager {
	cm := &ClientManager{clients: make(map[string]*ClientConnection), updateLeaderboardSig: make(chan struct{}, 1)}

	go cm.leaderboadDebounce()

	return cm
}

func (cm *ClientManager) ConnectClient(w http.ResponseWriter, r *http.Request, userID string, petService *service.PetService) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS ERROR] Upgrade failed for userID '%s': %v", userID, err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &ClientConnection{
		userID: userID,
		conn:   conn,
		cancel: cancel,
	}

	cm.mu.Lock()
	cm.clients[userID] = client
	activeCount := len(cm.clients)
	cm.mu.Unlock()

	log.Printf("[WS] Client connected successfully: userID='%s'. Total active clients: %d", userID, activeCount)

	go cm.startUpdates(ctx, userID, client, petService)
	go cm.listenAndServe(userID, client)
}

func (cm *ClientManager) startUpdates(ctx context.Context, userID string, client *ClientConnection, petService *service.PetService) {
	ticker := time.NewTicker(2 * time.Minute)
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

			event := map[string]any{
				"event_type": "pet.updated",
				"payload":    petResponse,
			}

			if err := client.WriteJSON(event); err != nil {
				return
			}

		case <-ctx.Done():
			log.Printf("[WS TICKER] Context done, stopping loop for userID: '%s'", userID)
			return
		}

	}
}

func (cm *ClientManager) listenAndServe(userID string, client *ClientConnection) {
	defer func() {
		client.cancel()

		cm.mu.Lock()
		delete(cm.clients, userID)
		cm.mu.Unlock()

		client.Close()
		log.Printf("[WS] Client '%s' disconnected", userID)
	}()

	for {
		if _, _, err := client.ReadMessage(); err != nil {
			return
		}
	}
}

func (cm *ClientManager) BroadcastLeaderboard() {
	log.Println("[WS BROADCAST] BroadcastLeaderboardUpdate called from service")
	select {
	case cm.updateLeaderboardSig <- struct{}{}:
		log.Println("[WS BROADCAST] Signal pushed to updateLeaderboardSig channel")
	default:
		log.Println("[WS BROADCAST] Channel full, signal skipped")
	}
}

func (cm *ClientManager) leaderboadDebounce() {
	log.Println("[WS DEBOUNCER] Debouncer started")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for range cm.updateLeaderboardSig {
		log.Println("[WS DEBOUNCER] Received signal from channel, waiting ticker")
		<-ticker.C

		event := map[string]any{
			"event_type": "leaderboard.position_updated",
		}

		cm.mu.RLock()
		clients := make([]*ClientConnection, 0, len(cm.clients))
		for _, client := range cm.clients {
			clients = append(clients, client)
		}
		cm.mu.RUnlock()

		for _, client := range clients {
			err := client.WriteJSON(event)
			if err != nil {
				log.Printf("[WS DEBOUNCER ERROR] Failed to send to userID '%s': %v", client.userID, err)
			}
		}
	}
}

func (cm *ClientManager) SendToClient(userID string, eventType string, v any) {
	log.Println("[WS CLIENT] sender started")

	event := map[string]any{
		"event_type": eventType,
		"payload":    v,
	}

	cm.mu.RLock()
	client, exists := cm.clients[userID]
	cm.mu.RUnlock()
	if !exists {
		log.Printf("[WS CLIENT] User '%s' is offline, event skipped", userID)
		return
	}

	if err := client.WriteJSON(event); err != nil {
		log.Printf("[WS CLIENT ERROR] Failed to send to userID '%s': %v", userID, err)
	}
}
