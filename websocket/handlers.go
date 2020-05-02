package websocket

import (
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/app/ptr"
	"github.com/peterzernia/app/restaurant"
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

func (h *Hub) handleJoin(message Message, conn *websocket.Conn) ([]restaurant.Restaurant, *Party) {
	if id, ok := message.Payload["id"].(string); ok {
		for _, party := range h.parties {
			ID, err := strconv.ParseInt(id, 10, 64)
			if err == nil && *party.ID == ID {
				conns := append(party.Conns, conn)
				party.Conns = conns

				options := restaurant.Options{
					Latitude:  ptr.Float64(52.492495),
					Longitude: ptr.Float64(13.393264),
					Limit:     ptr.Int64(50),
					Offset:    ptr.Int64(0),
					Radius:    ptr.Float64(1000),
				}
				search, err := restaurant.HandleList(options)

				if err == nil {
					party.Remaining = ptr.Int64(*search.Total - int64(len(search.Businesses)))
					return search.Businesses, &party
				}
			}
		}
	}

	return nil, nil
}
