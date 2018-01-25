package main

import (
	"amazonReviewCralwer/crawler"
	"amazonReviewCralwer/models"
	"amazonReviewCralwer/util"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var products []models.Product

func main() {

	for {
		if ur, err := getTopReviewUrl(); err == nil {
			if users := crawler.CrawlerTopReviewUser(ur.Url, util.Country(ur.Country)); users != nil {
				sendUsers(users)
				for _, u := range users {
					util.SetCountry(util.Country(u.Country))
					ps := crawler.GetReviewListUrl(u)
					for _, p := range ps {
						products = append(products, p)
						if len(products) >= 20 {
							crawler.SendProduct(products)
							products = make([]models.Product, 0)
						}
					}
				}
			}

		} else {
			util.Logger.Error(err.Error())
		}
	}

}

func sendUsers(users []models.User) {
	if len(users) <= 0 {
		return
	}

	j, err := json.Marshal(&users)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}

	body := bytes.NewBuffer(j)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://45.76.220.102:1323/user", body)
	if err != nil {
		util.Logger.Error(err.Error())
	}
	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Content-Encoding", "gzip")
	// Fetch Request
	_, err = client.Do(req)

	if err != nil {
		util.Logger.Error(err.Error())
	}
}

func getTopReviewUrl() (*UserJson, error) {
	// 获取URl (GET http://127.0.0.1:1323/url)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "http://45.76.220.102:1323/url", nil)
	if err != nil {
		return nil, err
	}
	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userJson UserJson
	err = json.Unmarshal(respBody, &userJson)
	if err != nil {
		return nil, err
	}
	return &userJson, nil
}

type UserJson struct {
	Url     string `json:"url"`
	Country int    `json:"country"`
}
