package websocket

import "github.com/peterzernia/lets-fork/ptr"

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

				// remove conn from party &
				// and remove party when no remaining connections
				// [TODO]: handle 1 remaining user
				if client.partyID != nil {
					for i, party := range h.parties {
						if *party.ID == *client.partyID {
							var index int
							for i, conn := range party.Conns {
								if client.conn == conn {
									index = i
									break
								}
							}
							h.parties[i].Conns[index] = party.Conns[len(party.Conns)-1]
							h.parties[i].Conns = party.Conns[:len(party.Conns)-1]

							if len(party.Conns) == 0 {
								h.parties[i] = h.parties[len(h.parties)-1]
								h.parties = h.parties[:len(h.parties)-1]
							}
							if len(h.parties[index].Conns) == 1 {
								h.parties[index].Status = ptr.String("waiting")
							}
						}
					}
				}
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

							// remove conn from party &
							// and remove party when no remaining connections
							// [TODO]: handle 1 remaining user
							if client.partyID != nil {
								for i, party := range h.parties {
									if *party.ID == *client.partyID {
										var index int
										for j, conn := range party.Conns {
											if client.conn == conn {
												index = j
												break
											}
										}
										h.parties[i].Conns[index] = party.Conns[len(party.Conns)-1]
										h.parties[i].Conns = party.Conns[:len(party.Conns)-1]

										if len(party.Conns) == 0 {
											h.parties[i] = h.parties[len(h.parties)-1]
											h.parties = h.parties[:len(h.parties)-1]
										}
										if len(h.parties[index].Conns) == 1 {
											h.parties[index].Status = ptr.String("waiting")
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
