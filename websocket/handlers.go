package websocket

import (
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/app/ptr"
	"github.com/peterzernia/app/restaurant"
)

func (h *Hub) handleCreate(c *Client) *Party {
	party := Party{
		ID:    ptr.Int64(int64(len(h.parties)) + 1),
		Conns: []*websocket.Conn{c.conn},
	}
	party.Likes = make(map[*websocket.Conn][]string)
	party.Likes[c.conn] = []string{}

	party.Matches = []restaurant.Restaurant{}

	c.partyID = party.ID

	parties := append(h.parties, party)
	h.parties = parties

	return &party
}

func (h *Hub) handleJoin(message Message, c *Client) ([]restaurant.Restaurant, *Party) {
	if id, ok := message.Payload["party_id"].(string); ok {
		for i, party := range h.parties {
			ID, err := strconv.ParseInt(id, 10, 64)

			if err == nil && *party.ID == ID {
				c.partyID = party.ID
				conns := append(party.Conns, c.conn)
				h.parties[i].Conns = conns
				h.parties[i].Likes[c.conn] = []string{}

				options := restaurant.Options{
					Latitude:  ptr.Float64(52.492495),
					Longitude: ptr.Float64(13.393264),
					Limit:     ptr.Int64(50),
					Offset:    ptr.Int64(0),
					Radius:    ptr.Float64(1000),
				}
				search, err := restaurant.HandleList(options)

				if err == nil {
					h.parties[i].Remaining = ptr.Int64(*search.Total - int64(len(search.Businesses)))
					h.parties[i].Current = search.Businesses
					h.parties[i].Restaurants = search.Businesses
					return search.Businesses, &h.parties[i]
				}
			}
		}
	}

	return nil, nil
}

func (h *Hub) handleSwipRight(message Message, c *Client) *Party {
	for i, party := range h.parties {
		if *party.ID == *c.partyID {
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
