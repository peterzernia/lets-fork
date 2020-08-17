package websocket

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/lets-fork/ptr"
	"github.com/peterzernia/lets-fork/restaurant"
	"github.com/peterzernia/lets-fork/utils"
	"github.com/stretchr/testify/assert"
)

// Check matches first with 2 users, then with 3 users
func TestCheckMatches(t *testing.T) {
	require := assert.New(t)

	_, err := utils.InitRDB()
	require.NoError(err)

	party := Party{
		Restaurants: []restaurant.Restaurant{
			{ID: ptr.String("R1")},
			{ID: ptr.String("R2")},
			{ID: ptr.String("R3")},
		},
	}

	user := User{
		ID:    ptr.String("U1"),
		Likes: []string{"R1"},
	}

	// checkMatches accesses user like info from the
	// redis db, so each user must be saved to db first
	err = setUser(user)
	require.NoError(err)

	conn := &websocket.Conn{}
	client := Client{
		conn: conn,
		id:   user.ID,
	}

	userTwo := User{
		ID:    ptr.String("U2"),
		Likes: []string{},
	}

	err = setUser(userTwo)
	require.NoError(err)

	connTwo := &websocket.Conn{}
	clientTwo := Client{
		conn: connTwo,
		id:   userTwo.ID,
	}

	clis := []Client{client, clientTwo}

	// No matches initially
	matches := party.checkMatches(clis)
	require.Empty(matches)

	userTwo.Likes = []string{"R1", "R2", "R3"}

	err = setUser(userTwo)
	require.NoError(err)

	matches = party.checkMatches(clis)
	require.Equal(len(matches), 1)
	require.Equal(*matches[0].ID, "R1")

	userThree := User{
		ID:    ptr.String("U3"),
		Likes: []string{},
	}

	err = setUser(userThree)
	require.NoError(err)

	connThree := &websocket.Conn{}
	clientThree := Client{
		conn: connThree,
		id:   userThree.ID,
	}

	clis = []Client{client, clientTwo, clientThree}
	matches = party.checkMatches(clis)
	require.Empty(matches)

	userThree.Likes = []string{"R1", "R2", "R3"}

	err = setUser(userThree)
	require.NoError(err)

	matches = party.checkMatches(clis)
	require.Equal(len(matches), 1)
	require.Equal(*matches[0].ID, "R1")
}
