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

const Cookie = "session-id=133-8500956-2009409; " +
	"session-id-time=2082787201l;" +
	" ubid-main=133-6893772-5740726;" +
	" s_cc=true; s_vnum=1947399226604&vn=1;" +
	" s_sq=[[B]]; s_nr=1515399305583-New; " +
	"s_dslv=1515399305584; s_ppv=100;" +
	" a-ogbcbff=1; " +
	"session-token=\"cZSlPMRt8u360s/ufDXaODK8+015vU0hYdlcvdE7Q95nAyjN8k+/2+GOwXFaQstzBtvCMEQAqoLREuME6kihWyFmD0WkzCL7UtAx81U4A3xAstdQtNC8HPa1/jSqtd14RY+eSpu695Lv2VHugVAo+n8qJBlOHAnhqDzSsIJxAyUtFmgGFssmVm7kByuuwil2tIr06Dq0CTFGh2MmBUf1eQ6oIhrVFFmKeFeEP09WlWK8HJzazZj5Kbm57t3l1pxiw6+Jhj5KM9tVsGXscLHY/A==\"; x-main=\"ezZ0qDyjzxneRSCR0IKIci1dA?DOIGkpwjYi0?TFgTbo5Dg72O8@dTNP4IZ@n4my\"; at-main=Atza|IwEBIHT7pByiU9me1KMaFTiD0r4OLIbVik0guwYh3mnGxvLelj2UjJDB-thUzAad5hI62iEKmxSGCWc0taQolNVTfzkhb2bxwHd3L6ldO4Acr2wVZ5cF4IyfOrWeiSVN0-rTo7eldGJR3ufwoFp5mspeKzxOFruM1JZx9x68aRnt3HTk6zpERVQHMpiQnktffDPjs3Yb2sFc0V-lwX1BtSbcZ2uWepwuOl7skiwIlCcIPVjaVVoq_6a6_mgvW2YE2BcZ65mUVPWk9QM3fyQcH17h-fTv-BTBonY_avjCCVbWZaVMMl2hmV4uklrdsog5r7zr8O0QVN1cKqVILo26E8WL0eLkWI6fCLatUn3t86XazjpRZux_mqlWHpAEhBP0cHyX9_CywtbqEt5rQjqIQQxSrn-v; sess-at-main=\"GP+xmj/5F+DSiAzermlfqyR5CwvOl2Zrfn7Iwd8uAtc=\"; sst-main=Sst1|PQHODdLpmKLUPsmWTUs5mpH4CFY1XAjRFaDM8Zly3R6fGRG1PXN-RR-BragU1OorONv41QnmEHIgx3WFk9QoEvyf2564ywr4WADsTp50fWbvQEOAQEnPK8JAX5EeIObF3lQCmJSX-1WK30c1Nj7KYwsrwSpubcvEKl6zw2mRke_DWGaNsJW1ImwLQh79V0V_JuK8B8hQx7SlM8d1EY9r1773HqQlurzRE2ZoHYv3RZ33QJVW8ycvqp9zZNzuDN6AFg1rC68jo08dst4_tGZ5CX57fg; lc-main=en_US; x-wl-uid=1bjXOXOzkydlFDdbs+hsDH18XuaxumxxDhvrejiwa7fbc9AC0HCkALwwstLnbInpMjOZwkf6ojugkuvXHP0yKP25KmErlFwTxwEpUXrOsvdu0o6MhOxn+sIEs56r6sucX5GY7CD2/y7w=; csm-hit=s-KYRJW9N1Y95HRP7BHNJ5|1515403111471"

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
	if err, s := getProfileHtml(user.ProfileUrl); err != nil {
		util.Logger.Error(err.Error())
	} else {
		if helpfulVotes, reviews, err := gethelpfulVotes(user.ProfileId); err == nil {
			user.HelpfulVotes = helpfulVotes
			user.Reviews = reviews
		} else {
			util.Logger.Error(err.Error())
		}
		user.Twitter = subStr(*s, "https://twitter.com/")
		user.Name = subStr(*s, "nameHeaderData")
		user.Facebook = subStr(*s, "https://www.facebook.com/")
		user.Instagram = subStr(*s, "https://www.instagram.com/")
		user.Youtube = subStr(*s, "https://www.youtube.com/")
		user.Pinterest = subStr(*s, "https://www.pinterest.com/")
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

	// Read Response Body
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

func getProfileHtml(profileUrl string) (err error, htmlstr *string) {

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
		return err, nil
	}

	defer resp.Body.Close()
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return err, nil
		}
	default:
		reader = resp.Body
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return err, nil
	}
	s := string(buf.Bytes())
	return nil, &s
}

func subStr(str, subStr string) string {
	if subStr == "nameHeaderData" {
		index := strings.Index(str, subStr)
		if index > 0 {
			tempStr := util.Substr(str, index+23, index+len(subStr)+100)
			if strings.Index(tempStr, "\"") > 0 {
				return util.Substr(tempStr, 0, strings.Index(tempStr, "\""))
			}
		}
	} else {
		index := strings.Index(str, subStr)
		if index > 0 {
			tempStr := util.Substr(str, index-2, index+len(subStr)+100)
			if strings.Index(tempStr, "\"") > 0 {
				return util.Substr(tempStr, 0, strings.Index(tempStr, "\""))
			}
		}
	}

	return ""
}

func gethelpfulVotes(userId string) (int, int, error) {

	client := &http.Client{}
	url := fmt.Sprintf("https://www.amazon.com/hz/gamification/api/contributor/dashboard/%s?ownerView=false&customerFollowEnabled=false", userId)
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
