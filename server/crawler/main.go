package main

import (
	"amazonReviewCralwer/crawler"
	"amazonReviewCralwer/models"
	"amazonReviewCralwer/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func main() {
	for {
		if err, users := getuser(); err == nil {
			for _, u := range users {
				util.SetCountry(util.Country(u.Country))
				crawler.GetReviewListUrl(u)
			}
		} else {
			util.Logger.Error(err.Error())
		}
	}
}

func getuser() (error, []models.User) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://45.76.220.102:1323/", nil)
	if err != nil {
		return err, nil
	}
	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return err, nil
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var users []models.User
	err = json.Unmarshal(respBody, &users)
	if err != nil {
		return err, nil
	}

	return nil, users
}
