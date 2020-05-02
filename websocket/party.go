package websocket

import "github.com/gorilla/websocket"

// Party represents 2+ users
type Party struct {
	ID    *int64            `json:"id"`
	Conns []*websocket.Conn `json:"-"`
}
