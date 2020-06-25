package restaurant

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/peterzernia/lets-fork/utils"
)

func handleGet(c *gin.Context) {
	var restaurant Restaurant
	rdb := utils.GetRDB()
	id := c.Param("id")

	rest, err := rdb.Get("restaurant:" + id).Result()
	if err == nil {
		err := json.Unmarshal([]byte(rest), &restaurant)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		// Sometimes yelp sends bad /GET requests with nil values
		// that should not be cached
		if restaurant.Name != nil {
			c.JSON(http.StatusOK, restaurant)
			return
		}
	}
	if err != redis.Nil {
		log.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/"+id, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_TOKEN"))
	res, err := client.Do(req)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// return json from GET request
	err = json.Unmarshal(body, &restaurant)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Cache restaurant for 24 hours in redis
	err = rdb.Set("restaurant:"+id, string(body), time.Hour*24).Err()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, restaurant)
}

// HandleList ...
func HandleList(options Options) (*SearchResponse, error) {
	search := SearchResponse{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_TOKEN"))

	// Add query parameters
	q := req.URL.Query()

	// Default to restaurants
	if options.Categories != nil {
		q.Add("categories", *options.Categories)
	} else {
		q.Add("categories", "restaurants")
	}
	q.Add("open_now", "true")
	q.Add("latitude", strconv.FormatFloat(*options.Latitude, 'f', -1, 64))
	q.Add("longitude", strconv.FormatFloat(*options.Longitude, 'f', -1, 64))
	q.Add("limit", strconv.FormatInt(*options.Limit, 10))
	q.Add("offset", strconv.FormatInt(*options.Offset, 10))
	q.Add("radius", strconv.FormatFloat(*options.Radius, 'f', -1, 64))
	if options.Price != nil {
		price, _ := json.Marshal(options.Price)
		q.Add("price", strings.Trim(string(price), "[]"))
	}

	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}
