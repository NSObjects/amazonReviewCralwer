package models

import "github.com/astaxie/beego/orm"

type Product struct {
	Id         int64
	Url        string
	CategoryId int64
	Categorys  []string `orm:"-" json:"categorys"`
	Name       string
	UserId     int64
}

func init() {
	orm.RegisterModel(new(Product))
}
