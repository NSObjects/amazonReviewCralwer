package main

import (
	"amazonReviewCralwer/crawler"
	"amazonReviewCralwer/models"
	"amazonReviewCralwer/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var products []models.Product

func main() {
	for {
		err, users := getuser()
		if err == nil {
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
		util.Logger.Error(string(respBody))
		return err, nil
	}

	return nil, users
}
