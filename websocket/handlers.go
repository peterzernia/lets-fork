package websocket

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/ptr"
	"github.com/peterzernia/lets-fork/restaurant"
)

// golang cannot unmarshal json into float fields
type options struct {
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
	Radius    string  `json:"radius"`
	Price     []int64 `json:"price"`
}

func (h *Hub) handleCreate(message Message, c *Client) *Party {
	party := Party{
		Conns: []*websocket.Conn{c.conn},
	}
	id, _ := h.generatePartyID()
	party.ID = &id

	party.Likes = make(map[*websocket.Conn][]string)
	party.Likes[c.conn] = []string{}

	party.Matches = []restaurant.Restaurant{}

	// Set options with workaround for *float64 fields
	options := options{}
	j, err := json.Marshal(message.Payload)
	if err != nil {
		log.Println(err, 1)
	}

	err = json.Unmarshal(j, &options)
	if err != nil {
		log.Println(err, 2)
	}

	lat, _ := strconv.ParseFloat(options.Latitude, 64)
	long, _ := strconv.ParseFloat(options.Longitude, 64)
	rad, _ := strconv.ParseFloat(options.Radius, 64)

	party.Options = &restaurant.Options{
		Latitude:  ptr.Float64(lat),
		Longitude: ptr.Float64(long),
		Limit:     ptr.Int64(50),
		Offset:    ptr.Int64(0),
		Radius:    ptr.Float64(rad),
		Price:     options.Price,
	}

	party.Status = ptr.String("waiting")
	c.partyID = party.ID

	parties := append(h.parties, party)
	h.parties = parties

	return &party
}

func (h *Hub) handleJoin(message Message, c *Client) (*Party, []*websocket.Conn) {
	if id, ok := message.Payload["party_id"].(string); ok {
		for i, party := range h.parties {
			if party.ID != nil && *party.ID == id {
				c.partyID = party.ID
				conns := append(party.Conns, c.conn)
				h.parties[i].Conns = conns
				h.parties[i].Likes[c.conn] = []string{}

				if len(party.Conns) == 1 { // Only the host is in the party
					search, err := restaurant.HandleList(*party.Options)

					if err == nil {
						h.parties[i].Current = search.Businesses
						h.parties[i].Total = search.Total
						h.parties[i].Restaurants = search.Businesses
						h.parties[i].Status = ptr.String("active")
						return &h.parties[i], h.parties[i].Conns
					}
				} else {
					// Reset matches when new user joins
					h.parties[i].Matches = []restaurant.Restaurant{}

					// Only send to new user
					return &h.parties[i], []*websocket.Conn{c.conn}
				}
			}
		}
	}

	return nil, nil
}

func (h *Hub) handleSwipRight(message Message, c *Client) *Party {
	for i, party := range h.parties {
		if c.partyID != nil && *party.ID == *c.partyID {
			if id, ok := message.Payload["restaurant_id"].(string); ok {
				h.parties[i].Likes[c.conn] = append(party.Likes[c.conn], id)
			}

			matches := party.checkMatches()
			if matches != nil {
				h.parties[i].Matches = matches
				return &h.parties[i]
			}
		}
	}

	return nil
}

func (h *Hub) handleRequestMore(message Message, c *Client) *Party {
	for i, party := range h.parties {
		if c.partyID != nil && *party.ID == *c.partyID {
			// Fetch more restaurants when they have not all been fetched
			if *party.Total-int64(len(party.Restaurants)) > 0 {
				h.parties[i].Options.Offset = ptr.Int64(int64(len(party.Restaurants)))
				search, err := restaurant.HandleList(*party.Options)

				if err == nil {
					h.parties[i].Current = search.Businesses
					h.parties[i].Restaurants = append(party.Restaurants, search.Businesses...)
					return &h.parties[i]
				}
			}
		}
	}

	return nil
}

func (h *Hub) handleQuit(c *Client) *Party {
	var index int
	var jndex int
	if c.partyID != nil {
		for i, party := range h.parties {
			if c.partyID != nil && *party.ID == *c.partyID {
				index = i
				for j, conn := range party.Conns {
					if c.conn == conn {
						jndex = j
					}
				}
			}
		}
		c.partyID = nil

		// Remove connection from party
		h.parties[index].Conns[jndex] = h.parties[index].Conns[len(h.parties[index].Conns)-1]
		h.parties[index].Conns = h.parties[index].Conns[:len(h.parties[index].Conns)-1]

		// Remove party if no connection
		if len(h.parties[index].Conns) == 0 {
			h.parties[index] = h.parties[len(h.parties)-1]
			h.parties = h.parties[:len(h.parties)-1]
			return nil
		}

		if len(h.parties[index].Conns) == 1 {
			h.parties[index].Status = ptr.String("waiting")
		}

		return &h.parties[index]
	}

	return nil
}
