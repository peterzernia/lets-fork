package websocket

import (
	"log"

	"github.com/peterzernia/lets-fork/restaurant"
)

// Party represents 2+ users
type Party struct {
	ID *string `json:"id,omitempty"`

	// Current 'batch' of fetched restaurants
	// to be added to the clients list of restaurants
	Current []restaurant.Restaurant `json:"current,omitempty"`
	Error   *string                 `json:"error,omitempty"`

	// Keep track of which restauraunts have been swiped right on by conn
	Matches     []restaurant.Restaurant `json:"matches,omitempty"`
	Options     *restaurant.Options     `json:"options"`
	Restaurants []restaurant.Restaurant `json:"restaurants,omitempty"`
	Status      *string                 `json:"status,omitempty"`
	Total       *int64                  `json:"total,omitempty"`
}

// Checks if any restaurant is liked by all the users
func (p *Party) checkMatches(clients []Client) []restaurant.Restaurant {
	matches := []restaurant.Restaurant{}
	counts := make(map[string]int)
	for _, c := range clients {
		user, err := getUser(*c.id)
		if err != nil {
			log.Println(err)
		}

		for _, restaurantID := range user.Likes {
			counts[restaurantID]++
		}
	}

	for restaurantID, count := range counts {
		if len(clients) != 1 && count == len(clients) {
			exists := false
			for _, match := range matches {
				if restaurantID == *match.ID {
					exists = true
				}
			}

			if !exists {
				for _, restaurant := range p.Restaurants {
					if restaurantID == *restaurant.ID {
						matches = append(matches, restaurant)
					}
				}
			}
		}
	}
	return matches
}
