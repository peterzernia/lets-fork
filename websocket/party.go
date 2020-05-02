package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/app/restaurant"
)

// Party represents 2+ users
type Party struct {
	ID      *int64                  `json:"id"`
	Conns   []*websocket.Conn       `json:"-"`
	Current []restaurant.Restaurant `json:"-"`

	// Keep track of which restauraunts have been swiped right on by conn
	Likes map[*websocket.Conn][]string `json:"-"`

	Matches []restaurant.Restaurant `json:"matches,omitempty"`

	// Remaining restaurants that have not been fetched
	Remaining   *int64                  `json:"-"`
	Restaurants []restaurant.Restaurant `json:"-"`
}

// Checks if any restaurant is liked by all the users
func (p *Party) checkMatches() []restaurant.Restaurant {
	counts := make(map[string]int)
	for _, conn := range p.Conns {
		for _, restaurantID := range p.Likes[conn] {
			counts[restaurantID]++
		}
	}
	fmt.Println(counts)
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
