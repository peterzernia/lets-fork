package websocket

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/ptr"
	"github.com/peterzernia/lets-fork/restaurant"
	"github.com/peterzernia/lets-fork/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetConnections(t *testing.T) {
	require := assert.New(t)

	hub := NewHub()
	go hub.Run()

	id := "000000"

	conn := &websocket.Conn{}
	client := Client{
		conn:    conn,
		id:      ptr.String("1"),
		partyID: ptr.String(id),
	}
	hub.register <- &client

	connTwo := &websocket.Conn{}
	clientTwo := Client{
		conn:    connTwo,
		id:      ptr.String("2"),
		partyID: ptr.String(id),
	}
	hub.register <- &clientTwo

	connThree := &websocket.Conn{}
	clientThree := Client{
		conn:    connThree,
		id:      ptr.String("3"),
		partyID: ptr.String("000001"),
	}
	hub.register <- &clientThree

	conns := hub.getConnections(id)
	require.Equal(len(conns), 2)
	require.Contains(conns, conn)
	require.Contains(conns, connTwo)
}

func TestGeneratePartyID(t *testing.T) {
	require := assert.New(t)

	_, err := utils.InitRDB()

	hub := NewHub()

	id, err := hub.generatePartyID()
	require.NoError(err)

	id2, err := hub.generatePartyID()
	require.NoError(err)

	require.NotEqual(id, id2)
}

func TestShuffle(t *testing.T) {
	require := assert.New(t)

	hub := NewHub()

	restaurants := []restaurant.Restaurant{
		{ID: ptr.String("1")},
		{ID: ptr.String("2")},
		{ID: ptr.String("3")},
	}

	shuffled := hub.shuffle(restaurants)
	require.Equal(len(shuffled), len(restaurants))
}
