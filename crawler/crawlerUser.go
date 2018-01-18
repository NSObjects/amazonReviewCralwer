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

	"sync"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/orm"
)

var Cookie = ""

type Country int

const (
	JAPAN Country = iota
	US
)

var countryURL = map[Country]string{
	JAPAN: "https://www.amazon.co.jp",
	US:    "https://www.amazon.com",
}

var (
	EmailNotFound  = errors.New("email not found")
	TokenNotFound  = errors.New("token not found")
	LastReviewPage = errors.New("last review paeg")
)

var BaseUrl string

func CrawlerTopReviewUser(c Country) {

	BaseUrl = countryURL[c]

	o := orm.NewOrm()
	for index := 1; index < 1000; index++ {
		err, q := getDocument(BaseUrl, index)
		if err != nil {
			util.Logger.Error(err.Error())
			continue
		}
		p := util.New(30)
		var wg sync.WaitGroup

		for _, user := range getUsers(q) {
			wg.Add(1)
			go func(u models.User) {
				p.Run(func() {
					if email, err := getUserEmail(u.ProfileUrl); err == nil {
						if email != "hidden@hidden.hidden" {
							u.Email = email
							if _, _, err := o.ReadOrCreate(&u, "profile_id"); err == nil {
								u.Country = int(c)
								configUser(&u)
								if _, err := o.Update(&u); err != nil {
									util.Logger.Error(err.Error())
								}
							} else {
								util.Logger.Error(err.Error())
							}
						}
					} else {
						if err != EmailNotFound {
							util.Logger.Error(err.Error())
						}
					}

					wg.Done()

				})

			}(user)
		}
		wg.Wait()
		p.Shutdown()
	}
}

func getUsers(q *goquery.Document) []models.User {

	var users []models.User
	var user models.User

	q.Find(".a-text-center").Find("a").Each(func(i int, selection *goquery.Selection) {
		if profileUrl, exit := selection.Attr("href"); exit {
			if strings.Contains(profileUrl, "account") {
				user.ProfileUrl = BaseUrl + util.Substr(profileUrl, 0, strings.Index(profileUrl, "/ref"))
			}
		}
		if name, exit := selection.Attr("name"); exit {
			user.ProfileId = name
			if user.ProfileId != "" && user.ProfileUrl != "" {
				users = append(users, user)
				user = models.User{}
			}
		}
	})

	return users
}

func configUser(user *models.User) {
	var userInfo struct {
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

	if err, s := getProfileHtml(user.ProfileUrl); err != nil {
		util.Logger.Error(err.Error())
	} else {
		if helpfulVotes, reviews, err := gethelpfulVotes(user.ProfileId); err == nil {
			user.HelpfulVotes = helpfulVotes
			user.Reviews = reviews
		} else {
			util.Logger.Error(err.Error())
		}

		d := strings.Index(s, "window.CustomerProfileRootProps = ")

		if d <= 0 {
			return
		}

		dd := s[d+len("window.CustomerProfileRootProps = ") : d+strings.Index(s[d:], "};")+1]
		err := json.Unmarshal([]byte(dd), &userInfo)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}

		user.Name = userInfo.NameHeaderData.Name
		user.Email = userInfo.BioData.PublicEmail
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
	}

}

func getDocument(url string, page int) (err error, g *goquery.Document) {

	client := &http.Client{}
	u := fmt.Sprintf("%s/hz/leaderboard/top-reviewers/ref=cm_cr_tr_link_%d?page=%d", url, page, page)
	req, err := http.NewRequest("GET", u, nil)

	// Headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", Cookie)
	req.Header.Add("Referer", u)

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
	req.Header.Add("Cookie", Cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Referer", profileUrl)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var emailJSON struct {
		Status string `json:"status"`
		Data   struct {
			Email string `json:"email"`
		} `json:"data"`
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(respBody, &emailJSON)
		if err != nil {
			return "", err
		} else {

			if emailJSON.Status == "ok" {
				return emailJSON.Data.Email, nil
			}
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
	req.Header.Add("Cookie", Cookie)
	req.Header.Add("If-None-Match", "W/\"4faafffbc5625b34d37ce3693b471d9c-gzip\"")
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

func subStr(str, subStr string) string {
	if subStr == "nameHeaderData" {
		index := strings.Index(str, subStr)
		if index > 0 && len(str) > index+len(subStr)+20 {
			tempStr := util.Substr(str, index+23, index+len(subStr)+20)
			if strings.Index(tempStr, "\"") > 0 {
				return util.Substr(tempStr, 0, strings.Index(tempStr, "\""))
			}
		}
	} else {
		index := strings.Index(str, subStr)
		if index > 0 && len(str) > index+len(subStr)+30 {
			tempStr := util.Substr(str, index-2, index+len(subStr)+30)
			if strings.Index(tempStr, "\"") > 0 {
				return util.Substr(tempStr, 0, strings.Index(tempStr, "\""))
			}
		}
	}

	return ""
}

func gethelpfulVotes(userId string) (int, int, error) {

	client := &http.Client{}
	url := fmt.Sprintf(BaseUrl+"/hz/gamification/api/contributor/dashboard/%s?ownerView=false&customerFollowEnabled=false", userId)
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	// Headers
	req.Header.Add("Cookie", Cookie)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return 0, 0, err
	}

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
