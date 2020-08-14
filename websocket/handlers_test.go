package websocket

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/ptr"
	"github.com/peterzernia/lets-fork/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandlersIntegration(t *testing.T) {
	// Test handleCreate
	require := assert.New(t)

	_, err := utils.InitRDB()
	require.NoError(err)

	hub := NewHub()
	go hub.Run()

	conn := &websocket.Conn{}

	payload := make(map[string]interface{})
	payload["categories"] = "restaurants"
	payload["latitude"] = "52.52"
	payload["longitude"] = "13.40"
	payload["radius"] = "5000"

	message := Message{
		Type:    "create",
		Payload: payload,
	}

	client := Client{
		conn: conn,
		id:   ptr.String("1"),
	}
	hub.register <- &client

	party, conns := hub.handleCreate(message, &client)
	require.NotNil(party.ID)
	require.Equal(*party.Status, "waiting")
	require.Equal(conns[0], conn)

	//////////////////////////////////////////////////
	// Test handleJoin
	connTwo := &websocket.Conn{}
	clientTwo := Client{
		conn: connTwo,
		id:   ptr.String("2"),
	}
	hub.register <- &clientTwo

	message.Type = "join"
	payload = make(map[string]interface{})
	payload["party_id"] = *party.ID
	message.Payload = payload

	party, conns = hub.handleJoin(message, &clientTwo)
	require.Equal(*party.Status, "active")
	require.Equal(len(conns), 2)
	require.Contains(conns, conn)
	require.Contains(conns, connTwo)

	//////////////////////////////////////////////////
	// Test handleSwipeRight
	require.Empty(party.Matches)

	id := party.Restaurants[0].ID
	message.Type = "swipe-right"
	payload = make(map[string]interface{})
	payload["restaurant_id"] = *id
	message.Payload = payload

	// initially nothing is returned as it is not a match yet
	party, conns = hub.handleSwipeRight(message, &client)
	require.Nil(party)
	require.Nil(conns)

	party, conns = hub.handleSwipeRight(message, &clientTwo)
	require.NotEmpty(party.Matches)
	require.Equal(*party.Matches[0].ID, *id)
	require.NotNil(conns)

	//////////////////////////////////////////////////
	// Test handleRequestMore
	restaurants := party.Restaurants
	message.Type = "request-more"
	message.Payload = nil

	party, conns = hub.handleRequestMore(message, &client)
	require.Greater(len(party.Restaurants), len(restaurants))
	require.NotNil(conns)

	//////////////////////////////////////////////////
	// Test handleQuit
	party, conns = hub.handleQuit(&clientTwo)
	require.Equal(len(conns), 1)
	require.Nil(clientTwo.partyID)
	require.Equal(*party.Status, "waiting")

	// rejoining should cause the matches to be cleared
	message.Type = "join"
	payload = make(map[string]interface{})
	payload["party_id"] = *party.ID
	message.Payload = payload

	party, conns = hub.handleJoin(message, &clientTwo)
	require.Equal(*party.Status, "active")
	require.Equal(len(conns), 2)
	require.Contains(conns, conn)
	require.Contains(conns, connTwo)
	require.NotNil(clientTwo.partyID)
	require.Equal(len(party.Matches), 0)

	// swiping right with clientTwo should cause
	// a match because the first clients swipe is saved
	message.Type = "swipe-right"
	payload = make(map[string]interface{})
	payload["restaurant_id"] = *id
	message.Payload = payload

	party, conns = hub.handleSwipeRight(message, &clientTwo)
	require.NotEmpty(party.Matches)
	require.Equal(*party.Matches[0].ID, *id)
	require.NotNil(conns)
}
