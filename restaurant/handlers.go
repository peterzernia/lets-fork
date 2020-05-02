package restaurant

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleGet(c *gin.Context) {
	id := c.Param("id")

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
	var resp interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
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
	q.Add("categories", "restaurants")
	q.Add("latitude", strconv.FormatFloat(*options.Latitude, 'f', -1, 64))
	q.Add("longitude", strconv.FormatFloat(*options.Longitude, 'f', -1, 64))
	q.Add("limit", strconv.FormatInt(*options.Limit, 10))
	q.Add("offset", strconv.FormatInt(*options.Offset, 10))
	q.Add("radius", strconv.FormatFloat(*options.Radius, 'f', -1, 64))

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
