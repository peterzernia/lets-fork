package websocket

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/peterzernia/lets-fork/utils"
)

func setParty(p Party) error {
	rdb := utils.GetRDB()

	jsn, err := json.Marshal(p)
	if err != nil {
		return err
	}

	rdb.Set("party:"+*p.ID, string(jsn), time.Hour)
	if err != nil {
		return err
	}
	return nil
}

func getParty(id string) (*Party, error) {
	var party Party
	rdb := utils.GetRDB()

	p, err := rdb.Get("party:" + id).Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}

	err = json.Unmarshal([]byte(p), &party)
	return &party, err
}
