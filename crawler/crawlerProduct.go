package crawler

import (
	"amazonReviewCralwer/models"
	"bytes"
	"fmt"
	"net/http"
	"unicode"

	"strings"

	"amazonReviewCralwer/util"

	"encoding/json"

	"io/ioutil"

	"sync"

	"github.com/PuerkitoBio/goquery"
)

func GetReviewListUrl(user models.User, baseUrl string) {
	var urls []string
	if token, reviewList, err := getReviewListToken(user.ProfileId, baseUrl); err == nil {
		for _, u := range reviewList {
			urls = append(urls, u)
		}
		for {
			if token, reviewList, err = getReviewList(token, baseUrl); err == nil {
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
				if productLink, err := getProductLint(reviewUrl); err == nil {
					link := baseUrl + productLink
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
						if len(products) > 10 {
							sendProduct(products)
							products = make([]models.Product, 0)
						}
					}

				} else {
					util.Logger.Error(err.Error())
				}
			})
		}(reviewUrl)

	}
	wg.Wait()
	p.Shutdown()

}

func sendProduct(products []models.Product) {
	if len(products) <= 0 {
		return
	}

	j, err := json.Marshal(&products)
	if err != nil {
		fmt.Println(err)
		return
	}

	body := bytes.NewBuffer(j)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "45.76.220.102:1323", body)
	if err != nil {
		fmt.Println(err)
	}
	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Content-Encoding", "gzip")
	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))

}

func getProductCategory(q *goquery.Document) (categoryList []string, err error) {

	q.Find(".a-link-normal.a-color-tertiary").Each(func(i int, selection *goquery.Selection) {
		s := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, selection.Text())

		if s != "" && s != "Reportabuse" {
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
	req.Header.Add("Cookie", Cookie)
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

func getReviewListToken(profileId string, baseUrl string) (token string, reviewList []string, err error) {

	client := &http.Client{}
	url := fmt.Sprintf(baseUrl+"/glimpse/timeline/%s?isWidgetOnly=true", profileId)
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	// Headers
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Accept", "text/html, */*; q=0.01")
	req.Header.Add("Referer", fmt.Sprintf(baseUrl+"/gp/profile/amzn1.account.%s", profileId))
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Add("Cookie", Cookie)
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

func getReviewList(token, baseUrl string) (nextToken string, reviewList []string, err error) {

	client := &http.Client{}
	url := fmt.Sprintf(baseUrl+"/glimpse/stories/next/ref=glimp_time_pag?token=%s&context=GlimpseTimeline&id&preview=false&dogfood=false", token)

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

func getProductLint(reviewUrl string) (productUrl string, err error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", reviewUrl, nil)
	if err != nil {
		return productUrl, err
	}
	// Headers
	req.Header.Add("Cookie", Cookie)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return productUrl, err
	}

	q, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return productUrl, err
	}

	var str string
	q.Find(".a-link-normal").Each(func(i int, selection *goquery.Selection) {
		str, _ = selection.Attr("data-hook")

		if str == "product-link" {
			if url, exist := selection.Attr("href"); exist {
				if url != "" {
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
