package main

import (
	"fmt"
	"os"

	"time"

	"io/ioutil"

	"strings"

	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/html"
)

//select user.id,count(*) from user,product where user.id=product.user_id Group by user.id order by count(*) Desc;
//#select count(*) from product where user_id=3
//#select category_id,COUNT(*) from product where user_id="2" and category_id > 0 GROUP BY category_id ORDER BY  COUNT(*) desc;
//#select * from product where user_id="9" and category_id = 476
//var s = " \"nameHeaderData\": {\"name\": \"粗茶ですが( ^-^)_旦‾\",\"profileExists\": true,\"inlineEditData\": null,\"isVerified\": false,\"urls\": {\"editButtonImageUrl\": \"//d1k8kvpjaf8geh.cloudfront.net/gp/profile/assets/icon_edit-0d9b7d9307686accef07de74ec135cb0c9847bd4a0cd810eeccb730723bc5b5c.png\" } },"

func main() {

	b, _ := ioutil.ReadFile("h.htm")
	s := string(b)

	d := strings.Index(s, "window.CustomerProfileRootProps = ")

	if d > 0 {
		dd := s[d+len("window.CustomerProfileRootProps = ") : d+strings.Index(s[d:], "};")+1]
		json.Unmarshal([]byte(dd), &userInfo)

	}
	//crawler.CrawlerTopReviewUser(crawler.US)
	//crawler.CrawlerProduct(crawler.US)
	//q := loadDoc()
	//fmt.Println(q.Find("a.a-link-normal").Attr("href"))
	//q.Find("a.a-link-normal").Each(func(i int, selection *goquery.Selection) {
	//	if s, exits := selection.Attr("data-hook"); exits {
	//		fmt.Println(s)
	//		fmt.Println(selection.Attr("href"))
	//	}
	//
	//})

}

func loadDoc() *goquery.Document {
	var f *os.File
	var e error

	if f, e = os.Open("h.htm"); e != nil {
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
	err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/amazon?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai", 30, 30)
	if err != nil {
		fmt.Println(err)
	}
}
