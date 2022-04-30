package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id         int
	Name       string
	Password   string
	Sex        int
	Mobile     string
	Email      string
	Address    string
	Status     int
	CreateTime int64
	UpdateTime int64
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User))
}

func NewUser() *User {
	return &User{}
}

// TableName 获取对应数据库表名.
func (u *User) TableName() string {
	return "user"
}

func (u *User) Add() error {
	o := orm.NewOrm()
	h := md5.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
	_, err := o.Insert(&u)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (u *User) GetList(pageIndex, pageSize int) ([]*User, error) {
	o := orm.NewOrm()
	offset := (pageIndex - 1) * pageSize
	var users []*User
	_, err := o.QueryTable(u.TableName()).Limit(pageSize, offset).All(&users)
	return users, err
}
func (u *User) TotalCount() (total int64, err error) {
	o := orm.NewOrm()
	total, err = o.QueryTable(u.TableName()).Count()
	return
}

func (u *User) GetUserInfoById(id int) (*User, error) {
	o := orm.NewOrm()
	var user User
	err := o.QueryTable(u.TableName()).Filter("id", id).One(&user)
	return &user, err
}

func (u *User) UserUpdate(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(u.TableName()).Filter("id", id).Update(orm.Params{
		"name":     u.Name,
		"email":    u.Email,
		"password": u.Password,
	})
	return
}

func (u *User) UserDel(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(u.TableName()).Filter("id", id).Delete()
	return
}

func (u *User) GetUserInfoByName(username string) (*User, error) {
	o := orm.NewOrm()
	var user User
	err := o.QueryTable(u.TableName()).Filter("name", username).One(&user)
	return &user, err
}
