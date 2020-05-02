package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/peterzernia/app/restaurant"
)

// Party represents 2+ users
type Party struct {
	ID    *int64            `json:"id"`
	Conns []*websocket.Conn `json:"-"`

	// Current 'batch' of fetched restaurants
	// to be added to the clients list of restaurants
	Current []restaurant.Restaurant `json:"current,omitempty"`

	// Keep track of which restauraunts have been swiped right on by conn
	Likes       map[*websocket.Conn][]string `json:"-"`
	Matches     []restaurant.Restaurant      `json:"matches,omitempty"`
	Options     *restaurant.Options          `json:"-"`
	Restaurants []restaurant.Restaurant      `json:"restaurants,omitempty"`
	Status      *string                      `json:"status"`
	Total       *int64                       `json:"total,omitempty"`
}

// Checks if any restaurant is liked by all the users
func (p *Party) checkMatches() []restaurant.Restaurant {
	counts := make(map[string]int)
	for _, conn := range p.Conns {
		for _, restaurantID := range p.Likes[conn] {
			counts[restaurantID]++
		}
	}

	for restaurantID, count := range counts {
		if count == len(p.Conns) {
			exists := false
			for _, match := range p.Matches {
				if restaurantID == *match.ID {
					exists = true
				}
			}

			if !exists {
				for _, restaurant := range p.Restaurants {
					if restaurantID == *restaurant.ID {
						p.Matches = append(p.Matches, restaurant)
						return p.Matches
					}
				}
			}
		}
	}
	return nil
}
