package routers

import (
	"github.com/beego/beego/v2/server/web"
	"liang_blog/controllers"
)

func init() {

	ns := web.NewNamespace("/admin",
		web.NSRouter("/index", &controllers.AdminController{}, "get:Index"),
		web.NSRouter("/welcome", &controllers.AdminController{}, "get:Welcome"),
		web.NSRouter("/member-list", &controllers.AdminController{}, "get:MemberList"),
		web.NSRouter("/member-add", &controllers.AdminController{}, "get,post:MemberAdd"),
		web.NSRouter("/member-edit", &controllers.AdminController{}, "get,post:MemberEdit"),
		web.NSRouter("/member-del", &controllers.AdminController{}, "get,post:MemberDel"),
		web.NSRouter("/login", &controllers.AdminController{}, "get,post:Login"),
		web.NSRouter("/out", &controllers.AdminController{}, "get,post:Out"),

		web.NSRouter("/article-list", &controllers.AdminController{}, "get:ArticleList"),
		web.NSRouter("/article-add", &controllers.AdminController{}, "get,post:ArticleAdd"),
		//web.NSRouter("/article-edit", &controllers.AdminController{}, "get,post:ArticleEdit"),
		//web.NSRouter("/article-del", &controllers.AdminController{}, "get,post:ArticleDel"),
	)
	//注册 namespace
	web.AddNamespace(ns)

	ns1 := web.NewNamespace("/index",
		web.NSRouter("/index", &controllers.IndexController{}, "get:Index"),
		web.NSRouter("/info", &controllers.IndexController{}, "get:Info"),
	)
	//注册 namespace
	web.AddNamespace(ns1)
}
