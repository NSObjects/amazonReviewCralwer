package main

import (
	"amazonReviewCralwer/crawler"
	"amazonReviewCralwer/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Country int

const (
	JAPAN Country = iota
	US
)

var countryURL = map[Country]string{
	JAPAN: "https://www.amazon.co.jp",
	US:    "https://www.amazon.com",
}

func main() {
	for {
		if err, users := getuser(); err == nil {
			for _, u := range users {
				crawler.GetReviewListUrl(u, countryURL[Country(u.Country)])
			}
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
