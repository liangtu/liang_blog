package models

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type Article struct {
	Id         int
	Name       string
	Content    string
	CreateTime int64
	UpdateTime int64
}

func NewArticle() *Article {
	return &Article{}
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Article))
}

func (a *Article) TableName() string {
	return "article"
}

func (a *Article) GetList(pageIndex, pageSize int) ([]*Article, error) {
	o := orm.NewOrm()
	offset := (pageIndex - 1) * pageSize
	var article []*Article
	_, err := o.QueryTable(a.TableName).Limit(pageSize, offset).All(&article)
	return article, err
}

func (a *Article) TotalCount() (total int64, err error) {
	o := orm.NewOrm()
	total, err = o.QueryTable(a.TableName).Count()
	return
}

func (a *Article) Add() error {
	o := orm.NewOrm()
	_, err := o.Insert(&a)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (a *Article) GetArticleInfoById(id int) (*Article, error) {
	o := orm.NewOrm()
	var article Article
	err := o.QueryTable(a.TableName).Filter("id", id).One(&article)
	return &article, err
}
