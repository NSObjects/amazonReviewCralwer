package main

import (
	"fmt"
	"os"

	"time"

	"amazonReviewCralwer/crawler"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/html"
)

//#select category_id,COUNT(*) from product where user_id="2" and category_id > 0 GROUP BY category_id ORDER BY  COUNT(*) desc;
//#select * from product where user_id="9" and category_id = 476
func main() {
	//crawler.CrawlerTopReviewUser(crawler.US)
	crawler.CrawlerProduct()
}

func loadDoc() *goquery.Document {
	var f *os.File
	var e error

	if f, e = os.Open("product.htm"); e != nil {
		panic(e.Error())
	}
	defer f.Close()

	var node *html.Node
	if node, e = html.Parse(f); e != nil {
		panic(e.Error())
	}

	return goquery.NewDocumentFromNode(node)
}

func init() {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		fmt.Println(err)
	}
	time.Local = local
	err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/amazon?charset=utf8&parseTime=true&loc=Asia%2FShanghai", 30, 30)
	if err != nil {
		fmt.Println(err)
	}
}
