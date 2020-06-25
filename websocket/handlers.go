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
	Categories *string `json:"categories"`
	Latitude   string  `json:"latitude"`
	Longitude  string  `json:"longitude"`
	Radius     string  `json:"radius"`
	Price      []int64 `json:"price"`
}

func (h *Hub) handleCreate(message Message, c *Client) (*Party, []*websocket.Conn) {
	party := Party{}

	id, _ := h.generatePartyID()
	party.ID = &id

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
		Categories: options.Categories,
		Latitude:   ptr.Float64(lat),
		Longitude:  ptr.Float64(long),
		Limit:      ptr.Int64(50),
		Offset:     ptr.Int64(0),
		Radius:     ptr.Float64(rad),
		Price:      options.Price,
	}

	party.Status = ptr.String("waiting")
	c.partyID = party.ID

	err = setParty(party)
	if err != nil {
		log.Println(err)
	}

	user := User{
		ID:      c.id,
		Likes:   []string{},
		PartyID: party.ID,
	}

	err = setUser(user)
	if err != nil {
		log.Println(err)
	}

	return &party, []*websocket.Conn{c.conn}
}

func (h *Hub) handleJoin(message Message, c *Client) (*Party, []*websocket.Conn) {
	if id, ok := message.Payload["party_id"].(string); ok {
		party, err := getParty(id)
		if err != nil {
			return nil, nil
		}

		c.partyID = ptr.String(id)

		conns := h.getConnections(id)

		user := User{
			ID:      c.id,
			Likes:   []string{},
			PartyID: c.partyID,
		}
		err = setUser(user)
		if err != nil {
			log.Println(err)
		}

		if party.Total == nil {
			search, err := restaurant.HandleList(*party.Options)

			if err == nil {
				party.Current = h.shuffle(search.Businesses)
				party.Total = search.Total
				party.Restaurants = search.Businesses
				party.Status = ptr.String("active")

				err = setParty(*party)
				if err != nil {
					log.Println(err)
				}

				return party, conns
			}
		} else {
			// Reset matches when new user joins
			party.Matches = []restaurant.Restaurant{}
			if len(conns) == 1 {
				party.Status = ptr.String("waiting")

				err = setParty(*party)
				if err != nil {
					log.Println(err)
				}

				return party, []*websocket.Conn{c.conn}
			}

			party.Status = ptr.String("active")

			err = setParty(*party)
			if err != nil {
				log.Println(err)
			}

			return party, conns
		}
	}

	return nil, nil
}

func (h *Hub) handleSwipeRight(message Message, c *Client) (*Party, []*websocket.Conn) {
	if c.partyID != nil {
		party, err := getParty(*c.partyID)
		if err != nil {
			log.Println(err)
		}

		user, err := getUser(*c.id)

		if id, ok := message.Payload["restaurant_id"].(string); ok {
			exists := false
			for _, restaurant := range user.Likes {
				if restaurant == id {
					exists = true
				}
			}
			if !exists {
				user.Likes = append(user.Likes, id)
			}
		}

		err = setUser(*user)
		if err != nil {
			log.Println(err)
		}

		conns := h.getConnections(*c.partyID)

		clients := []Client{}
		for cli := range h.clients {
			if cli.partyID != nil && *cli.partyID == *c.partyID {
				clients = append(clients, *cli)
			}
		}

		matches := party.checkMatches(clients)
		if matches != nil && len(matches) != len(party.Matches) {
			party.Matches = matches
			err = setParty(*party)
			if err != nil {
				log.Println(err)
			}

			return party, conns
		}
	}

	return nil, nil
}

func (h *Hub) handleRequestMore(message Message, c *Client) (*Party, []*websocket.Conn) {
	if c.partyID != nil {
		party, err := getParty(*c.partyID)
		if err != nil {
			log.Println(err)
		}

		// Fetch more restaurants when they have not all been fetched
		if party.Total != nil && *party.Total-int64(len(party.Restaurants)) > 0 {
			party.Options.Offset = ptr.Int64(int64(len(party.Restaurants)))
			search, err := restaurant.HandleList(*party.Options)

			if err == nil {
				current := h.shuffle(search.Businesses)
				party.Current = current
				party.Restaurants = append(party.Restaurants, current...)

				conns := h.getConnections(*c.partyID)

				err = setParty(*party)
				if err != nil {
					log.Println(err)
				}

				return party, conns
			}
		}
	}

	return nil, nil
}

func (h *Hub) handleQuit(c *Client) (*Party, []*websocket.Conn) {
	var id string
	if c.partyID != nil {
		id = *c.partyID
		c.partyID = nil
	}

	conns := h.getConnections(id)

	if len(conns) == 1 {
		party, err := getParty(id)
		if err != nil {
			log.Println(err)
		}

		if party == nil {
			return nil, nil
		}

		party.Status = ptr.String("waiting")
		err = setParty(*party)
		if err != nil {
			log.Println(err)
		}

		return party, conns
	}

	return nil, nil
}
