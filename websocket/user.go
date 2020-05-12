package websocket

// User stores a users info in case the server
// disconnects
type User struct {
	ID      *string  `json:"id"` // unique device id
	Likes   []string `json:"likes"`
	PartyID *string  `json:"party_id"`
}
