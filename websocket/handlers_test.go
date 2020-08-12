package websocket

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/ptr"
	"github.com/peterzernia/lets-fork/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreate(t *testing.T) {
	require := assert.New(t)

	_, err := utils.InitRDB()
	require.NoError(err)

	hub := NewHub()
	go hub.Run()

	conn := &websocket.Conn{}

	message := Message{
		Type:    "create",
		Payload: make(map[string]interface{}),
	}

	client := Client{
		conn: conn,
		id:   ptr.String("1"),
	}

	party, conns := hub.handleCreate(message, &client)
	require.NotNil(party.ID)
	require.Equal(*party.Status, "waiting")
	require.Equal(conns[0], conn)
}
