package crawler

import (
	"amazonReviewCralwer/models"
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"strings"

	"amazonReviewCralwer/util"

	"encoding/json"

	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

func GetReviewListUrl(user models.User) []models.Product {

	switch util.Country(user.Country) {
	case util.US:
		return us(user)
	case util.JAPAN:
		return jp(user)
	}

	return nil

}

func us(user models.User) []models.Product {
	var urls []string
	if token, reviewList, err := getReviewListToken(user.ProfileId); err == nil {
		for _, u := range reviewList {
			urls = append(urls, u)
		}
		for {
			if token, reviewList, err = getReviewList(token); err == nil {
				for _, u := range reviewList {
					urls = append(urls, u)
				}
			} else {
				if err != LastReviewPage {
					util.Logger.Error(err.Error())
				}
				break
			}
		}
	} else {
		if err != TokenNotFound {
			util.Logger.Error(err.Error())
		}
	}

	var products []models.Product
	p := util.New(30)
	var wg sync.WaitGroup

	for _, reviewUrl := range urls {
		wg.Add(1)
		go func(url string) {
			p.Run(func() {
				if productLink, err := getProductLint(url, util.Country(user.Country)); err == nil {
					link := util.BaseUrl + productLink
					product := models.Product{
						UserId: user.Id,
						Url:    link,
					}
					if doc, err := getProductDoc(link); err == nil {
						if categoryList, err := getProductCategory(doc); err == nil {
							product.Categorys = categoryList
						} else {
							util.Logger.Error(err.Error())
						}
						if name, exists := doc.Find("#imgTagWrapperId").Children().Attr("alt"); exists {
							product.Name = name
						} else if name = doc.Find("#ebooksProductTitle").Text(); name != "" {
							product.Name = name
						} else if name = doc.Find("#productTitle").Text(); name != "" {
							product.Name = name
						}
						products = append(products, product)
					}

				} else {
					util.Logger.Error(err.Error())
				}
			})
			wg.Done()
		}(reviewUrl)

	}
	wg.Wait()
	p.Shutdown()

	return products
}

func jp(user models.User) (products []models.Product) {
	work := util.New(30)
	var wg sync.WaitGroup
	url := user.ProfileUrl
	i := strings.Index(url, "account.") + 8
	reviews := getJPReviewList(url[i:])
	for _, p := range reviews {
		wg.Add(1)
		go func(r Reviews) {
			work.Run(func() {
				link := util.BaseUrl + r.Urls.ProductURL
				product := models.Product{
					UserId: user.Id,
					Url:    link,
					Name:   r.ProductTitle,
				}
				if doc, err := getProductDoc(link); err == nil {
					if categoryList, err := getProductCategory(doc); err == nil {
						product.Categorys = categoryList
					} else {
						util.Logger.Error(err.Error())
					}

					products = append(products, product)
				} else {
					util.Logger.Error(err.Error())
				}
			})
			wg.Done()
		}(p)

	}
	wg.Wait()
	work.Shutdown()

	return
}

type Reviews struct {
	ProductTitle string `json:"productTitle"`
	Urls         struct {
		ProductURL string `json:"productUrl"`
	} `json:"urls"`
}

func getJPReviewList(profileId string) (reviewLists []Reviews) {

	var index int

	for {
		var reviewList struct {
			Reviews    []Reviews `json:"reviews"`
			NextOffset int       `json:"nextOffset"`
		}
		client := &http.Client{}
		url := fmt.Sprintf("https://www.amazon.co.jp/gp/profile/amzn1.account.%s/activity_feed?review_offset=%d", profileId, index*10)

		req, err := http.NewRequest("GET", url, nil)

		// Headers
		req.Header.Add("Accept", "*/*")
		req.Header.Add("Cookie", "x-wl-uid=1aQWRn7TwQAv4mngfMGdQQMzqbVvaJypQekW0WHfP8uXS2AeNgCe5Wvf1L8eKSJ9k8HDZvBCqyVqlYhjyLmUUkB+nnLKYzfY8pjqHL4RL93iDe7jKSsXX2B/3qEEBA6xlGTAtubLojl0=; ubid-acbjp=357-4101823-5570864; session-token=fdHy8mZNRH4UIKML0PI/evaCgSFNTQaxBV58mazqKARB9NAxj91a6oO4SiOuTC54DbVcFrVghy/sy/8RDhg3uIPJSVPcRKwBP9iqs2lC8PG7UluAI9yk8HR8DeV70YJG7H8TXbbImF4asHvCKtMOlN1r0unmKrVT17FAINyy3mJx0XI00hrSwEE2cg4bgC7nkmH8jjs2ohTmZdnUzavoVCcPUInuHuLetBH/503wSWr/JJ2jEr+UNHG3TnOnhrRruG8M+Xgthuk=; session-id-time=2082726001l; session-id=355-9214000-6210429")
		req.Header.Add("Referer", "https://www.amazon.co.jp/gp/profile/amzn1.account.AGBZ2GB6CLXKYJSCPUCZLR72ZNYQ?language=en_US")
		req.Header.Add("Host", "www.amazon.co.jp")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
		req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		parseFormErr := req.ParseForm()
		if parseFormErr != nil {
			util.Logger.Error(err.Error())
			break
		}

		// Fetch Request
		resp, err := client.Do(req)

		if err != nil {
			util.Logger.Error(err.Error())
			break
		}

		// Read Response Body
		respBody, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(respBody, &reviewList)
		if err != nil {
			util.Logger.Error(url)
			util.Logger.Error(err.Error())
			break
		}

		if len(reviewList.Reviews) <= 0 {
			break
		}
		for _, v := range reviewList.Reviews {
			reviewLists = append(reviewLists, v)
		}

		index++

	}

	return

}

func SendProduct(products []models.Product) {
	if len(products) <= 0 {
		return
	}

	j, err := json.Marshal(&products)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}

	body := bytes.NewBuffer(j)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://45.76.220.102:1323/", body)
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

func getProductCategory(q *goquery.Document) (categoryList []string, err error) {

	q.Find("a.a-link-normal.a-color-tertiary").Each(func(i int, selection *goquery.Selection) {
		s := strings.TrimSpace(selection.Text())
		if s != "Report abuse" && s != "" {
			categoryList = append(categoryList, s)
		}
	})

	return categoryList, nil
}

func getProductDoc(url string) (q *goquery.Document, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Accept", "text/html, */*; q=0.01")
	req.Header.Add("Referer", url)
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", util.Cookie)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	q, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return
}

func getReviewListToken(profileId string) (token string, reviewList []string, err error) {

	client := &http.Client{}
	url := fmt.Sprintf(util.BaseUrl+"/glimpse/timeline/%s?isWidgetOnly=true", profileId)
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	// Headers
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Accept", "text/html, */*; q=0.01")
	req.Header.Add("Referer", fmt.Sprintf(util.BaseUrl+"/gp/profile/amzn1.account.%s", profileId))
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", util.Cookie)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return "", nil, err
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return "", nil, err
	}

	q, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil, err
	}

	token, _ = q.Find(".glimpse-main-pagination-trigger").Attr("data-pagination-token")
	if token == "" {
		return "", nil, TokenNotFound
	}

	q.Find(".a-link-normal").Each(func(i int, selection *goquery.Selection) {
		if href, exist := selection.Attr("href"); exist {
			if strings.Contains(href, "review") {
				reviewList = append(reviewList, href)
			}
		}
	})

	return
}

func getReviewList(token string) (nextToken string, reviewList []string, err error) {

	client := &http.Client{}
	url := fmt.Sprintf(util.BaseUrl+"/glimpse/stories/next/ref=glimp_time_pag?token=%s&context=GlimpseTimeline&id&preview=false&dogfood=false", token)

	req, err := http.NewRequest("GET", url, nil)

	// Headers
	//req.Header.Add("Cookie", "x-main=\"Z3iPGlEC2POisUHBDRzlkShywbR0m45lcpgZl6a6LwhCtHbgPs9@PGX4v?yOmDBb\"; x-wl-uid=1oZ0/FGCuddDdb/vZeyyd45aqSM/HtRaQxNrhrqbIxII6DfbpF/mxM0y/VFz7K/w+dMtzPjWgYTEGkJEtHV0ZIw==; ubid-main=133-6893772-5740726; session-id=133-8500956-2009409; session-id-time=2082787201l; session-token=\"IzlB3HOYVF6/7JlihCzTQgBZw9lUf9+YNTTFBJeKj7y0nUAL+KX3QrNg6vARxkW0Zd8aIMC/PVbG8jJIYy0oe7yXtLaaGqZ/Q8KjdCidV1LoZAyEU03HDrbAm1Gyk24Ckh4H2O/9Mw8/HXWNgbKahkkVmOn6zireixCfWX4NMvr2Y63TJ63uGM+vNRqv31kpRrDPycO3ogZ6E8e2i+KCVRHXNZU4Z6eiTyb7kYJj9JEBrD2GYAo/8ub+12lLS5wwGi8Tu4emxx+Jr8EauWG14w==\"")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return "", nil, err
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return "", nil, err
	}

	q, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil, err
	}

	q.Find(".a-link-normal").Each(func(i int, selection *goquery.Selection) {
		if href, exist := selection.Attr("href"); exist {
			if strings.Contains(href, "review") {
				reviewList = append(reviewList, href)
			}
		}
	})

	t, exit := q.Find(".glimpse-main-pagination-trigger").Attr("data-pagination-token")
	if exit == true && t != "" && t != token {
		return t, reviewList, nil
	}

	return "", nil, LastReviewPage
}

func getProductLint(reviewUrl string, c util.Country) (productUrl string, err error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", reviewUrl, nil)
	if err != nil {
		return productUrl, err
	}
	// Headers
	req.Header.Add("Cookie", util.Cookie)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return productUrl, err
	}
	defer resp.Body.Close()
	q, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return productUrl, err
	}

	q.Find("a.a-link-normal").Each(func(i int, selection *goquery.Selection) {
		if s, exits := selection.Attr("data-hook"); exits {
			if s == "product-link" {

				if url, exits := selection.Attr("href"); exits {
					//b, _ := ioutil.ReadAll(resp.Body)
					//fmt.Println(string(b))
					if c == util.JAPAN {
						url += "&language=en_US"
					}
					productUrl = url
				}
			}

		}
	})

	if productUrl != "" {
		return productUrl, nil
	}

	return "", fmt.Errorf("product url not found /n %s", reviewUrl)
}
