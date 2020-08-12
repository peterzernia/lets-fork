package restaurant

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/peterzernia/lets-fork/ptr"
	"github.com/peterzernia/lets-fork/utils"
	"github.com/stretchr/testify/assert"
)

// Testing that the yelp endpoints are still returning correctly
func TestHandleList(t *testing.T) {
	require := assert.New(t)

	limit := ptr.Int64(50)

	opts := Options{
		Categories: ptr.String("restaurants"),
		Latitude:   ptr.Float64(52.52),
		Longitude:  ptr.Float64(13.40),
		Limit:      limit,
		Offset:     ptr.Int64(0),
		Radius:     ptr.Float64(5000),
		Price:      []int64{1, 2, 3},
	}

	res, err := HandleList(opts)

	require.NoError(err)
	require.Equal(int64(len(res.Businesses)), *limit)
}

func TestHandleGet(t *testing.T) {
	require := assert.New(t)

	id := "0gxqSYCeVWilIOOOjZtSdA"

	rdb, err := utils.InitRDB()
	require.NoError(err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: id}}

	handleGet(c)

	require.Equal(w.Code, 200)

	// Restaurant should not be nil in the redis database initially
	_, err = rdb.Get("restaurant:" + id).Result()
	require.NotEqual(err, redis.Nil)
}
