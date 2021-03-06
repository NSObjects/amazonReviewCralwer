package main

import (
	"os"

	"time"

	"amazonReviewCralwer/util"

	"fmt"

	"strings"

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
	//s := "https://www.amazon.co.jp/gp/profile/amzn1.account.AEPFMU6MODEHFQ253YZMIOZAJ4CQ/ref=cm_cr_tr_tbl_511_name"
	//i := strings.Index(s, "account.") + 8
	//fmt.Println(s[i : i+strings.Index(s[i:], "/")])

	//crawler.CrawlerTopReviewUser(crawler.US)
	//crawler.CrawlerProduct(crawler.US)
	q := loadDoc()
	//fmt.Println(q.Find("a.a-link-normal").Attr("href"))
	q.Find("a.a-text-center").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Find("a.a-profile.avatar.a-spacing-none").Attr("href"))
		if v, exits := selection.Attr("name"); exits {
			fmt.Println(v)
		}
	})

	q.Find(".a-text-center").Find("a").Each(func(i int, selection *goquery.Selection) {
		if profileUrl, exit := selection.Attr("href"); exit {
			if strings.Contains(profileUrl, "account") {
				fmt.Println(profileUrl)
			}
		}
		if name, exit := selection.Attr("name"); exit {
			fmt.Println(name)
		}
	})
	//q.Find("a").Each()
	//q.Find(".a-profile.avatar.a-spacing-none").Each(func(i int, selection *goquery.Selection) {
	//	if s, exits := selection.Attr("data-hook"); exits {
	//		if s == "product-link" {
	//			fmt.Println(selection.Attr("href"))
	//		}
	//
	//	}
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
		util.Logger.Error(err.Error())
	}
	time.Local = local
	err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/amazon?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai", 30, 30)
	if err != nil {
		util.Logger.Error(err.Error())
	}
}
