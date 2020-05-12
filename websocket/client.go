package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/ptr"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	partyID *string
	id      *string
}

func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		message := Message{}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Printf("Websocket error: %s", err)
		}

		response := Response{}
		switch message.Type {
		case "create":
			party, conns := c.hub.handleCreate(message, c)
			res, err := json.Marshal(party)
			if err != nil {
				log.Println(err)
			}
			response.Res = res
			response.Conns = conns
		case "join":
			party, conns := c.hub.handleJoin(message, c)
			if conns != nil {
				res, err := json.Marshal(party)
				if err != nil {
					log.Println(err)
				}
				response.Res = res
				response.Conns = conns
			} else {
				party := Party{Error: ptr.String("Party does not exist")}
				res, err := json.Marshal(party)
				if err != nil {
					log.Println(err)
				}
				response.Res = res
				response.Conns = []*websocket.Conn{c.conn}
			}
		case "swipe-right":
			party, conns := c.hub.handleSwipeRight(message, c)
			if party != nil {
				res, err := json.Marshal(party)
				if err != nil {
					log.Println(err)
				}
				response.Res = res
				response.Conns = conns
			}
		case "request-more":
			party, conns := c.hub.handleRequestMore(message, c)
			if party != nil {
				res, err := json.Marshal(party)
				if err != nil {
					log.Println(err)
				}
				response.Res = res
				response.Conns = conns
			}
		case "quit":
			party, conns := c.hub.handleQuit(c)
			if party != nil {
				res, err := json.Marshal(party)
				if err != nil {
					log.Println(err)
				}
				response.Res = res
				response.Conns = conns
			}
		default:
			res := []byte("Unrecognized message type" + message.Type)
			log.Println("Unrecognized message type" + message.Type)

			response.Res = res
			response.Conns = []*websocket.Conn{c.conn}
		}

		c.hub.broadcast <- response
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Serve handles websocket requests from the peer.
func serve(hub *Hub, w http.ResponseWriter, r *http.Request) {
	var id string
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if len(r.URL.Query()["id"]) > 0 {
		id = r.URL.Query()["id"][0]
	}

	// in the case client was disconnected, user will have
	// been stored in the rdb
	user, err := getUser(id)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		hub:     hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		id:      ptr.String(id),
		partyID: nil,
	}

	// return user info from rdb
	if user != nil {
		client.partyID = user.PartyID
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
