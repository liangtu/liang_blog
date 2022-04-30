package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	username, _ := config.String("username")
	password, _ := config.String("password")
	host, _ := config.String("host")
	port, _ := config.String("port")
	database, _ := config.String("database")
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		fmt.Println(err)
	}
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, password, host, port, database) + "&loc=Asia%2FShanghai"
	orm.RegisterDataBase("default", "mysql", dbLink)
}
