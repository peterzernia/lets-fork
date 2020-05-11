package restaurant

// Options represents search options
type Options struct {
	Latitude  *float64
	Longitude *float64
	Limit     *int64
	Offset    *int64
	Radius    *float64
	Price     []int64
}

// SearchResponse represents a response from yelp's search endpoint
type SearchResponse struct {
	Businesses []Restaurant `json:"businesses"`
	Total      *int64       `json:"total"`
}

// Restaurant represents a restaurant
type Restaurant struct {
	Alias        *string      `json:"alias"`
	Categories   []Category   `json:"categories"`
	Coordinates  *Coordinates `json:"coordinates"`
	DisplayPhone *string      `json:"display_phone"`
	Hours        []Hours      `json:"hours"`
	ID           *string      `json:"id"`
	ImageURL     *string      `json:"image_url"`
	IsClaimed    *bool        `json:"is_claimed"`
	IsClosed     *bool        `json:"is_closed"`
	Location     *Location    `json:"location"`
	Name         *string      `json:"name"`
	Phone        *string      `json:"string"`
	Photos       []string     `json:"photos"`
	Price        *string      `json:"price"`
	Rating       *float64     `json:"rating"`
	Transactions []string     `json:"transactions"`
	URL          *string      `json:"url"`
}

// Category represents a restaurant category
type Category struct {
	Alias *string `json:"alias"`
	Title *string `json:"title"`
}

// Coordinates is the coordinates of the restaurant
type Coordinates struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// Location is the address of the restaurant
type Location struct {
	Address1 *string `json:"address1"`
	Address2 *string `json:"address2"`
	Address3 *string `json:"address3"`
	City     *string `json:"city"`
	Country  *string `json:"country"`
	State    *string `json:"state"`
	ZipCode  *string `json:"zip_code"`
}

// Hours represents the hours of a restaurant
type Hours struct {
	Open      []Hour  `json:"open"`
	HoursType *string `json:"hours_type"`
	IsOpenNow *bool   `json:"is_open_now"`
}

// Hour represents the open hours of a restaurant for one day of the week
type Hour struct {
	IsOvernight *bool   `json:"is_overnight"`
	Start       *string `json:"start"`
	End         *string `json:"end"`
	Day         *int64  `json:"day"`
}
