package websocket

import (
	"github.com/peterzernia/lets-fork/utils"
)

// Hub ...
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Response

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Keep track of all the live parties
	parties []Party
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Response),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		parties:    []Party{},
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
	for _, party := range h.parties {
		// Check to make sure id doesn't exist
		// and run generatePartyID again if it does
		if id == *party.ID {
			id, err = h.generatePartyID()
			if err != nil {
				return "", err
			}
		}
	}

	return id, nil
}
