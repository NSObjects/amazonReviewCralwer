package main

import (
	"amazonReviewCralwer/models"
	"amazonReviewCralwer/util"

	"net/http"
	"time"

	"github.com/labstack/echo"

	"encoding/json"

	"io/ioutil"

	"strings"

	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var countryURL = map[util.Country]string{
	util.JAPAN: "https://www.amazon.co.jp",
	util.US:    "https://www.amazon.com",
}

var (
	lastId      chan int64
	country     = util.US
	cruuentPage chan int
)

func main() {

	e := echo.New()

	lastId = make(chan int64, 1)
	lastId <- 0
	cruuentPage = make(chan int, 1)
	cruuentPage <- 0

	e.GET("/", func(context echo.Context) error {
		o := orm.NewOrm()
		var users []models.User
		l := <-lastId
		_, err := o.QueryTable("user").
			Filter("id__gt", l).
			OrderBy("id").
			Limit(10).
			All(&users)
		if err != nil {
			lastId <- 0
			return context.String(http.StatusBadRequest, err.Error())
		}

		if len(users) == 0 || l == users[len(users)-1].Id {
			lastId <- 0
			return context.String(http.StatusBadRequest, "last")
		}
		lastId <- users[len(users)-1].Id
		return context.JSON(http.StatusOK, users)
	})

	e.POST("/", func(context echo.Context) error {

		b, err := ioutil.ReadAll(context.Request().Body)
		if err != nil {
			util.Logger.Error(err.Error())
			return context.String(200, err.Error())
		}

		var products []models.Product
		if err = json.Unmarshal(b, &products); err != nil {
			util.Logger.Error(err.Error())
			return context.String(200, err.Error())
		}
		saveProducts(products)
		return context.String(200, "")
	})

	e.GET("/url", func(context echo.Context) error {
		var userJson struct {
			Url     string `json:"url"`
			Country int    `json:"country"`
		}
		url := topReviewUserUrl()
		userJson.Url = url
		userJson.Country = int(country)
		return context.JSON(http.StatusOK, userJson)
	})

	e.POST("/user", func(context echo.Context) error {
		b, err := ioutil.ReadAll(context.Request().Body)
		if err != nil {
			util.Logger.Error(err.Error())
			return context.String(200, err.Error())
		}

		var users []models.User
		if err = json.Unmarshal(b, &users); err != nil {
			util.Logger.Error(err.Error())
			return context.String(200, err.Error())
		}
		saveUsers(users)
		return context.String(200, "")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func saveUsers(user []models.User) {
	o := orm.NewOrm()
	for _, u := range user {
		if _, _, err := o.ReadOrCreate(&u, "profile_url"); err != nil {
			util.Logger.Error(err.Error())
		}
	}
}

func saveProducts(products []models.Product) {
	o := orm.NewOrm()
	for _, product := range products {
		u := models.User{
			ProfileUrl: product.UserProfile,
		}
		if err := o.Read(&u, "profile_url"); err == nil {
			product.UserId = u.Id
			if _, _, err := o.ReadOrCreate(&product, "user_id", "url"); err == nil {
				var categoryId int64 = 0
				for _, c := range product.Categorys {
					category := models.Category{
						Name: c,
					}
					if created, _, err := o.ReadOrCreate(&category, "name"); err == nil {
						if created {
							category.ParentId = categoryId
							if _, err := o.Update(&category, "parent_id"); err != nil {
								util.Logger.Error(err.Error())
							}
						}
						categoryId = category.Id
					} else {
						util.Logger.Error(err.Error())
					}
				}
				product.CategoryId = categoryId
				if _, err := o.Update(&product, "category_id", "name", "review_url"); err != nil {
					util.Logger.Error(err.Error())
				}
			} else {
				if strings.Contains(err.Error(), "Error 1062: Duplicate entry") == false {
					util.Logger.Error(err.Error())
				}
			}
		} else {
			util.Logger.Error(err.Error())
		}

	}

}

func topReviewUserUrl() string {
	page := <-cruuentPage
	url := countryURL[country]
	if page >= 1000 {
		cruuentPage <- 0
		if country == util.JAPAN {
			country = util.US
		} else if country == util.US {
			country = util.JAPAN
		}
	} else {
		cruuentPage <- page + 1
	}

	return fmt.Sprintf("%s/hz/leaderboard/top-reviewers/ref=cm_cr_tr_link_%d?page=%d", url, page, page)
}

func init() {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		util.Logger.Error(err.Error())
	}
	time.Local = local
	err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/amazon?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai", 30, 30)
	if err != nil {
		util.Logger.Error(err.Error())
	}
}
