package main

import (
	"github.com/beego/beego/v2/server/web"
	_ "liang_blog/models"
	_ "liang_blog/routers"
)

func main() {
	web.Run()
}
