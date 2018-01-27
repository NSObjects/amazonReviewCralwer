package crawler

import (
	"amazonReviewCralwer/models"
	"amazonReviewCralwer/util"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"

	"strings"

	"io/ioutil"

	"encoding/json"

	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type UserInfo struct {
	CustomerID     string `json:"customerId"`
	NameHeaderData struct {
		Name string `json:"name"`
	} `json:"nameHeaderData"`
	BioData struct {
		PublicEmail string `json:"publicEmail"`
		Social      struct {
			HasLinks    bool `json:"hasLinks"`
			SocialLinks []struct {
				Type string      `json:"type"`
				URL  interface{} `json:"url"`
			} `json:"socialLinks"`
		} `json:"social"`
	} `json:"bioData"`
}

var (
	EmailNotFound  = errors.New("email not found")
	TokenNotFound  = errors.New("token not found")
	LastReviewPage = errors.New("last review paeg")
)

func CrawlerTopReviewUser(url string, c util.Country) (users []models.User) {
	util.SetCountry(c)
	if err, q := getDocument(url); err == nil {
		return getUsers(q)
	}

	return nil
}

func configUser(user models.User, c util.Country) models.User {

	if userInfo, err := getUserInfo(user.ProfileUrl); err == nil {
		if helpfulVotes, reviews, err := gethelpfulVotes(user.ProfileId); err == nil {
			user.HelpfulVotes = helpfulVotes
			user.Reviews = reviews
		} else {
			util.Logger.Error(err.Error())
		}
		user.Email = userInfo.BioData.PublicEmail
		user.Name = userInfo.NameHeaderData.Name
		if userInfo.BioData.Social.HasLinks == true {
			for _, v := range userInfo.BioData.Social.SocialLinks {

				switch v.Type {
				case "facebook":
					if s, ok := v.URL.(string); ok {
						user.Facebook = s
					}
				case "twitter":
					if s, ok := v.URL.(string); ok {
						user.Twitter = s
					}
				case "pinterest":
					if s, ok := v.URL.(string); ok {
						user.Pinterest = s
					}
				case "instagram":
					if s, ok := v.URL.(string); ok {
						user.Instagram = s
					}
				case "youtube":
					if s, ok := v.URL.(string); ok {
						user.Youtube = s
					}

				}
			}
		}
		user.Country = int(c)
	} else if err != EmailNotFound {
		util.Logger.Error(err.Error())
	}

	return user
}

func getUsers(q *goquery.Document) []models.User {

	var users []models.User
	var user models.User

	q.Find(".a-text-center").Find("a").Each(func(i int, selection *goquery.Selection) {
		if profileUrl, exit := selection.Attr("href"); exit {
			if strings.Contains(profileUrl, "account") {
				user.ProfileUrl = util.BaseUrl + util.Substr(profileUrl, 0, strings.Index(profileUrl, "/ref"))
			}
		}
		if name, exit := selection.Attr("name"); exit {
			user.ProfileId = name
			if user.ProfileId != "" && user.ProfileUrl != "" {
				user = configUser(user, util.Country(user.Country))
				if user.Email != "" && user.Email != "hidden@hidden.hidden" {
					users = append(users, user)
				}

				user = models.User{}
			}
		}
	})

	return users
}

func getUserInfo(profileUrl string) (userInfo *UserInfo, err error) {

	if err, s := getProfileHtml(profileUrl); err == nil {

		d := strings.Index(s, "window.CustomerProfileRootProps = ")

		if d <= 0 {
			util.Logger.Error(profileUrl)
			return nil, errors.New("CustomerProfileRootProps not found")
		}

		dd := s[d+len("window.CustomerProfileRootProps = ") : d+strings.Index(s[d:], "};")+1]
		err = json.Unmarshal([]byte(dd), &userInfo)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	if userInfo.BioData.PublicEmail == "" || userInfo.BioData.PublicEmail == "hidden@hidden.hidden" {
		return nil, EmailNotFound
	}

	return
}

func getDocument(url string) (err error, g *goquery.Document) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	// Headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", util.Cookie)
	req.Header.Add("Referer", url)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return err, nil
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return err, nil
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			util.Logger.Error(err.Error())
		}
	default:
		reader = resp.Body
	}
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)

	g, err = goquery.NewDocumentFromReader(buf)
	if err != nil {
		return err, nil
	}
	return
}

func getUserEmail(profileUrl string) (email string, err error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", profileUrl+"/customer_email", nil)
	if err != nil {
		return "", err
	}
	// Headers
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", util.Cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Referer", profileUrl)
	req.Header.Add("accept-encoding", "gzip, deflate, br")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			util.Logger.Error(err.Error())
		}
	default:
		reader = resp.Body
	}
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)

	var emailJSON struct {
		Status string `json:"status"`
		Data   struct {
			Email string `json:"email"`
		} `json:"data"`
	}
	fmt.Println(string(buf.Bytes()))

	err = json.Unmarshal(buf.Bytes(), &emailJSON)
	if err != nil {
		return "", err
	} else {

		if emailJSON.Status == "ok" {
			return emailJSON.Data.Email, nil
		}
	}

	return "", EmailNotFound
}

func getProfileHtml(profileUrl string) (err error, htmlstr string) {

	client := &http.Client{}
	// Create request
	req, err := http.NewRequest("GET", profileUrl, nil)

	// Headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", util.Cookie)
	req.Header.Add("Referer", profileUrl)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return err, htmlstr
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return err, htmlstr
		}
	default:
		reader = resp.Body
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return err, htmlstr
	}
	s := string(buf.Bytes())
	return nil, s
}

func gethelpfulVotes(userId string) (int, int, error) {

	client := &http.Client{}
	url := fmt.Sprintf(util.BaseUrl+"/hz/gamification/api/contributor/dashboard/%s?ownerView=false&customerFollowEnabled=false", userId)
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	// Headers
	req.Header.Add("Cookie", util.Cookie)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	var helpfulVotes struct {
		HelpfulVotes struct {
			HelpfulVotesData struct {
				Count string `json:"count"`
			} `json:"helpfulVotesData"`
		} `json:"helpfulVotes"`
		Reviews struct {
			ReviewsCountData struct {
				Count string `json:"count"`
			} `json:"reviewsCountData"`
		} `json:"reviews"`
	}

	err = json.Unmarshal(respBody, &helpfulVotes)
	if err != nil {
		return 0, 0, err
	}
	var count int
	if strings.Contains(helpfulVotes.HelpfulVotes.HelpfulVotesData.Count, "k") {
		s := strings.Replace(helpfulVotes.HelpfulVotes.HelpfulVotesData.Count, "k", "", -1)

		c, err := strconv.ParseFloat(s, 64)
		if err == nil {
			count = int(c * 1000)
		}
	} else {

		c, _ := strconv.Atoi(helpfulVotes.HelpfulVotes.HelpfulVotesData.Count)
		count = c
	}

	var reviews int
	if strings.Contains(helpfulVotes.Reviews.ReviewsCountData.Count, "k") {
		s := strings.Replace(helpfulVotes.Reviews.ReviewsCountData.Count, "k", "", -1)

		c, err := strconv.ParseFloat(s, 64)
		if err == nil {
			reviews = int(c * 1000)
		}
	} else {

		c, _ := strconv.Atoi(helpfulVotes.HelpfulVotes.HelpfulVotesData.Count)
		reviews = c
	}

	return count, reviews, nil
}
