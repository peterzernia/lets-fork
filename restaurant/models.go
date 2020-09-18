package restaurant

// Options represents search options
type Options struct {
	Categories *string  `json:"categories,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Limit      *int64   `json:"limit,omitempty"`
	Offset     *int64   `json:"offset,omitempty"`
	OpenNow    *bool    `json:"open_now,omitempty"`
	Radius     *float64 `json:"radius,omitempty"`
	Price      []int64  `json:"price,omitempty"`
}

// SearchResponse represents a response from yelp's search endpoint
type SearchResponse struct {
	Businesses []Restaurant `json:"businesses"`
	Total      *int64       `json:"total"`
}

// Restaurant represents a restaurant
type Restaurant struct {
	Alias        *string      `json:"alias,omitempty"`
	Categories   []Category   `json:"categories,omitempty"`
	Coordinates  *Coordinates `json:"coordinates,omitempty"`
	DisplayPhone *string      `json:"display_phone,omitempty"`
	Hours        []Hours      `json:"hours,omitempty"`
	ID           *string      `json:"id,omitempty"`
	ImageURL     *string      `json:"image_url,omitempty"`
	IsClaimed    *bool        `json:"is_claimed,omitempty"`
	IsClosed     *bool        `json:"is_closed,omitempty"`
	Location     *Location    `json:"location,omitempty"`
	Name         *string      `json:"name,omitempty"`
	Phone        *string      `json:"string,omitempty"`
	Photos       []string     `json:"photos,omitempty"`
	Price        *string      `json:"price,omitempty"`
	Rating       *float64     `json:"rating,omitempty"`
	ReviewCount  *int64       `json:"review_count,omitempty"`
	Transactions []string     `json:"transactions,omitempty"`
	URL          *string      `json:"url,omitempty"`
}

// Category represents a restaurant category
type Category struct {
	Alias *string `json:"alias,omitempty"`
	Title *string `json:"title,omitempty"`
}

// Coordinates is the coordinates of the restaurant
type Coordinates struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// Location is the address of the restaurant
type Location struct {
	Address1 *string `json:"address1,omitempty"`
	Address2 *string `json:"address2,omitempty"`
	Address3 *string `json:"address3,omitempty"`
	City     *string `json:"city,omitempty"`
	Country  *string `json:"country,omitempty"`
	State    *string `json:"state,omitempty"`
	ZipCode  *string `json:"zip_code,omitempty"`
}

// Hours represents the hours of a restaurant
type Hours struct {
	Open      []Hour  `json:"open,omitempty"`
	HoursType *string `json:"hours_type,omitempty"`
	IsOpenNow *bool   `json:"is_open_now,omitempty"`
}

// Hour represents the open hours of a restaurant for one day of the week
type Hour struct {
	IsOvernight *bool   `json:"is_overnight,omitempty"`
	Start       *string `json:"start,omitempty"`
	End         *string `json:"end,omitempty"`
	Day         *int64  `json:"day,omitempty"`
}
