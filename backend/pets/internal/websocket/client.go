package websocket

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientConnection struct {
	userID string
	conn   *websocket.Conn
	mu     sync.Mutex
	cancel context.CancelFunc
}

func (cc *ClientConnection) WriteJSON(v any) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.conn.WriteJSON(v)
}

func (cc *ClientConnection) Close() error {
	return cc.conn.Close()
}

func (cc *ClientConnection) ReadMessage() (int, []byte, error) {
	return cc.conn.ReadMessage()
}
