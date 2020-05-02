package restaurant

// Options represents search options
type Options struct {
	Latitude  *float64
	Longitude *float64
	Limit     *int64
	Offset    *int64
	Radius    *float64
}

// SearchResponse represents a response from yelp's search endpoint
type SearchResponse struct {
	Businesses []Restaurant `json:"businesses"`
	Total      *int64       `json:"total"`
}

// Restaurant represents a restaurant
type Restaurant struct {
	Alias *string `json:"alias"`
	// Categories   []Categories `json:"categories"`
	// Coordinates  *Coordinates `json:"coordinates"`
	DisplayPhone *string `json:"display_phone"`
	// Hours        *Hours       `json:"hours"`
	ID        *string `json:"id"`
	ImageURL  *string `json:"image_url"`
	IsClaimed *bool   `json:"is_claimed"`
	IsClosed  *bool   `json:"is_closed"`
	// Location  *Location `json:"location"`
	Name         *string  `json:"name"`
	Phone        *string  `json:"string"`
	Photos       []string `json:"photos"`
	Price        *string  `json:"price"`
	Rating       *float64 `json:"rating"`
	Transactions []string `json:"transactions"`
	URL          *string  `json:"url"`
}
