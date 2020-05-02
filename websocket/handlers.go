package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/peterzernia/app/ptr"
)

func (h *Hub) handleCreate(conn *websocket.Conn) Party {
	party := Party{
		ID:    ptr.Int64(int64(len(h.parties)) + 1),
		Conns: []*websocket.Conn{conn},
	}

	for client := range h.clients {
		if client.conn == conn {
			client.partyID = party.ID
		}
	}

	parties := append(h.parties, party)
	h.parties = parties

	return party
}
