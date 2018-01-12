package crawler

import (
	"amazonReviewCralwer/models"
	"fmt"
	"net/http"
	"unicode"

	"strings"

	"amazonReviewCralwer/util"

	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/orm"
)

func CrawlerProduct(c Country) {
	baseUrl = countryURL[c]
	var users []models.User
	o := orm.NewOrm()

	_, err := o.QueryTable("user").Limit(10000).All(&users)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	p := util.New(30)
	var wg sync.WaitGroup
	for _, u := range users {
		wg.Add(1)
		go func(user models.User) {
			p.Run(func() {
				getReviewProduct(user)
				wg.Done()
			})
		}(u)
	}

	wg.Wait()
	p.Shutdown()
}

func getReviewProduct(user models.User) {

	token, reviewList, err := getReviewListToken(user.ProfileId)
	if err != nil || token == "" {
		if err != TokenNotFound {
			util.Logger.Error(err.Error())
		}
		return
	}

	inserProduct(reviewList, user.Id)

	for {
		token, reviewList, err = getReviewList(token)
		if err != nil || len(reviewList) == 0 {
			if err != LastReviewPage {
				util.Logger.Error(err.Error())
			}
			break
		}
		inserProduct(reviewList, user.Id)
	}
}

func inserProduct(reviewList []string, userId int64) {
	if userId == 0 {
		panic(userId)
	}
	o := orm.NewOrm()
	for _, reviewUrl := range reviewList {
		if productLink, err := getProductLint(reviewUrl); err == nil {
			link := baseUrl + productLink
			product := models.Product{
				UserId: userId,
				Url:    link,
			}

			if _, _, err := o.ReadOrCreate(&product, "user_id", "url"); err == nil {
				if categoryList, err := getProductCategory(link); err == nil {
					var categoryId int64 = 0
					for _, c := range categoryList {
						category := models.Category{
							Name: c,
						}
						if created, id, err := o.ReadOrCreate(&category, "name"); err == nil {
							if created {
								category.ParentId = categoryId
								if _, err := o.Update(&category, "parent_id"); err != nil {
									util.Logger.Error(err.Error())
								}
							}
							categoryId = id
						} else {
							util.Logger.Error(err.Error())
						}
					}
					product.CategoryId = categoryId
				} else {
					util.Logger.Error(err.Error())
				}
				if _, err := o.Update(&product, "category_id"); err != nil {
					util.Logger.Error(err.Error())
				}
			} else {
				util.Logger.Error(err.Error())
			}

		} else {
			util.Logger.Error(err.Error())
		}
	}
}

func getProductCategory(url string) (categoryList []string, err error) {

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

	q, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	//categoryList := make([]string, 0)
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

func getReviewListToken(profileId string) (token string, reviewList []string, err error) {

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
	if len(token) == 0 {
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
