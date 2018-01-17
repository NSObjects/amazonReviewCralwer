package main

import (
	"amazonReviewCralwer/models"
	"amazonReviewCralwer/util"
	"bytes"

	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"encoding/json"

	"compress/gzip"

	"amazonReviewCralwer/crawler"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	lastId chan int64
)

func main() {

	go func() {
		for {
			crawler.CrawlerTopReviewUser(crawler.JAPAN)
			crawler.CrawlerTopReviewUser(crawler.US)
		}
	}()
	e := echo.New()

	lastId = make(chan int64, 1)
	lastId <- 0

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
			return context.String(http.StatusBadRequest, err.Error())
		}
		lastId <- users[len(users)-1].Id
		return context.JSON(http.StatusOK, users)
	})

	e.POST("/", func(context echo.Context) error {

		reader, err := gzip.NewReader(context.Request().Body)
		if err != nil {
			fmt.Println(err)
		}

		var b []byte
		buf := bytes.NewBuffer(b)
		buf.ReadFrom(reader)
		var products []models.Product
		if err = json.Unmarshal(buf.Bytes(), &products); err != nil {
			fmt.Println(err)
		}
		saveProducts(products)
		return context.String(200, "")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func saveProducts(products []models.Product) {
	o := orm.NewOrm()
	for _, product := range products {
		if _, _, err := o.ReadOrCreate(&product, "user_id", "url"); err == nil {
			var categoryId int64 = 0
			for _, c := range product.Categorys {
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
			if _, err := o.Update(&product, "category_id", "name"); err != nil {
				util.Logger.Error(err.Error())
			}
		} else {
			util.Logger.Error(err.Error())
		}
	}

}

func init() {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		fmt.Println(err)
	}
	time.Local = local
	err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/amazon?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai", 30, 30)
	if err != nil {
		fmt.Println(err)
	}
}
