package websocket

import "github.com/gorilla/websocket"

// Message represents a message sent from the client
type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// Response is the marshalled JSON res, and
// the list of conns to send the response to
type Response struct {
	Conns []*websocket.Conn
	Res   []byte
}
