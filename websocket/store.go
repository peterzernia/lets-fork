package websocket

import (
	"encoding/json"
	"time"

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

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(p), &party)
	return &party, err
}

func setUser(u User) error {
	rdb := utils.GetRDB()

	jsn, err := json.Marshal(u)
	if err != nil {
		return err
	}

	rdb.Set("user:"+*u.ID, string(jsn), time.Hour)
	if err != nil {
		return err
	}
	return nil
}

func getUser(id string) (*User, error) {
	var user User
	rdb := utils.GetRDB()

	p, err := rdb.Get("user:" + id).Result()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(p), &user)
	return &user, err
}
