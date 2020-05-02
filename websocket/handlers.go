package websocket

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/app/ptr"
)

func (h *Hub) handleCreate(conn *websocket.Conn) *Party {
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

	return &party
}

func (h *Hub) handleJoin(message Message, conn *websocket.Conn) *Party {
	fmt.Println(message.Payload)
	if id, ok := message.Payload["id"].(string); ok {
		for _, party := range h.parties {
			ID, err := strconv.ParseInt(id, 10, 64)
			if err == nil && *party.ID == ID {
				conns := append(party.Conns, conn)
				party.Conns = conns
			}

			return &party
		}
	}

	return nil
}
