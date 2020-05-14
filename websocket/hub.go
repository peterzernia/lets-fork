package websocket

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/utils"
)

// Hub ...
type Hub struct {
	// Inbound messages from the clients
	broadcast chan Response

	// Registered clients
	clients map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Response),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run runs the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case response := <-h.broadcast:
			// Send the response to the clients in the response.Conns array
			for client := range h.clients {
				for _, conn := range response.Conns {
					if client.conn == conn {
						select {
						case client.send <- response.Res:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}

// Returns the list of connections associated with a party
func (h *Hub) getConnections(partyID string) []*websocket.Conn {
	ids := []string{}
	conns := []*websocket.Conn{}
	for cli := range h.clients {
		if cli.partyID != nil && *cli.partyID == partyID {
			// skip duplicates
			skip := false
			for _, id := range ids {
				if *cli.id == id {
					skip = true
				}
			}
			if !skip {
				conns = append(conns, cli.conn)
				ids = append(ids, *cli.id)
			}
		}
	}

	return conns
}

// generatePartyID returns a random 6 digit number as string.
func (h *Hub) generatePartyID() (string, error) {
	const letters = "0123456789"
	bytes, err := utils.GenerateRandomBytes(6)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	id := string(bytes)

	_, err = getParty(id)
	if err != redis.Nil {
		id, err = h.generatePartyID()
		if err != nil {
			return "", err
		}
	}

	return id, nil
}
