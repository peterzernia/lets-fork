package restaurant

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

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
