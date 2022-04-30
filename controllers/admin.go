package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"liang_blog/models"
	pagination "liang_blog/utiils"
	"math"
	"strings"
	"time"
)

type AdminController struct {
	web.Controller
}

type Response struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
}

type UserResData struct {
	Id         int
	Name       string
	Sex        string
	Mobile     string
	Email      string
	Address    string
	Status     string
	CreateTime string
}

type ArticleResData struct {
	Id         int
	Name       string
	Content    string
	CreateTime string
}

//验证是否登录
func (c *AdminController) Prepare() {
	uri := c.Ctx.Input.URI()
	//把不需要登录的绕开
	//判断是否相等

	v := c.GetSession("passwordMd5")
	fmt.Println(v)
	if v == nil && uri != "/admin/login" {
		urlss := c.URLFor("AdminController.Login")
		c.Redirect(urlss, 302)
	}
}

func (c *AdminController) URLMapping() {
	c.Mapping("Index", c.Index)
	c.Mapping("MemberList", c.MemberList)
	c.Mapping("MemberAdd", c.MemberAdd)
	c.Mapping("MemberDel", c.MemberDel)
	c.Mapping("ArticleList", c.ArticleList)
	c.Mapping("ArticleAdd", c.ArticleAdd)
	c.Mapping("Login", c.Login)
	c.Mapping("Out", c.Out)
}
func (c *AdminController) Login() {
	valid := validation.Validation{}
	if c.Ctx.Input.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")
		if v := valid.Required(username, "username").Message("用户名不能为空！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Required(password, "password1").Message("密码必填！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		//去数据库匹配是否有个这个用户 密码是否相等
		user := models.NewUser()
		userInfo, err := user.GetUserInfoByName(username)
		if err != nil {
			c.Json(300, "用户名或者密码错误")
			return
		}
		h := md5.New()
		h.Write([]byte(password))
		passwordMd5 := hex.EncodeToString(h.Sum(nil))
		if userInfo.Password != passwordMd5 {
			c.Json(300, "用户名或者密码错误")
			return
		}
		//有这个用户让其登录写入session
		v := c.GetSession("passwordMd5")
		if v == nil {
			c.SetSession("passwordMd5", int(1))
		} else {
			c.SetSession("passwordMd5", v.(int)+1)
		}

		c.Json(200, "登录成功")
	} else {
		c.TplName = "admin/login.html"
	}
}

func (c *AdminController) Out() {
	c.DelSession("passwordMd5")
	c.Redirect(c.URLFor("AdminController.Login"), 302)
}
func (c *AdminController) Index() {
	c.TplName = "admin/index.html"
}

func (c *AdminController) Welcome() {

	c.TplName = "admin/welcome.html"
}

func (c *AdminController) MemberList() {

	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 10

	user := models.NewUser()
	list, err := user.GetList(pageIndex, pageSize)
	if err != nil {

	}
	var userList []*UserResData
	for _, v := range list {
		SexNote := "男"
		if v.Sex == 2 {
			SexNote = "女"
		}
		StatusNote := "启用"
		if v.Status == 2 {
			StatusNote = "禁用"
		}
		userList = append(userList, &UserResData{
			Id:         v.Id,
			Name:       v.Name,
			Sex:        SexNote,
			Mobile:     v.Mobile,
			Email:      v.Email,
			Address:    v.Address,
			Status:     StatusNote,
			CreateTime: time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05"),
		})
	}

	c.Data["data"] = userList
	totalCount, _ := user.TotalCount()

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, int(totalCount), pageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	c.Data["totalCount"] = totalCount

	c.TplName = "admin/member-list.html"
}

func (c *AdminController) MemberAdd() {

	valid := validation.Validation{}

	if c.Ctx.Input.IsPost() {
		username := c.GetString("username")
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		email := c.GetString("email")

		if v := valid.Required(username, "username").Message("用户名不能为空！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Required(username, "password1").Message("密码必填！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Required(username, "password2").Message("密码必填！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Email(email, "email").Message("邮箱格式不对！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if password2 != password1 {
			c.Json(300, "两次密码不等")
			return
		}

		//添加数据库
		user := models.NewUser()
		user.Name = username
		user.Address = ""
		user.Email = email
		user.Status = 1
		user.Mobile = "18898998976"
		user.Password = password1
		user.CreateTime = time.Now().Unix()
		user.UpdateTime = time.Now().Unix()
		user.Sex = 1
		if err := user.Add(); err != nil {
			c.Json(300, err.Error())
			return
		}
		c.Json(200, "添加成功")
	} else {

		c.TplName = "admin/member-add.html"
	}
}

func (c *AdminController) MemberEdit() {
	if c.Ctx.Input.IsPost() {
		valid := validation.Validation{}
		id, _ := c.GetInt("id")
		username := c.GetString("username")
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		email := c.GetString("email")

		if v := valid.Required(username, "username").Message("用户名不能为空！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Required(username, "password1").Message("密码必填！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Required(username, "password2").Message("密码必填！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Email(email, "email").Message("邮箱格式不对！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if password2 != password1 {
			c.Json(300, "两次密码不等")
			return
		}

		user := models.NewUser()
		user.Name = username
		user.Email = email
		user.Password = password1
		if err := user.UserUpdate(id); err != nil {
			c.Json(300, err.Error())
			return
		}
		c.Json(200, "编辑成功")
	} else {
		id, _ := c.GetInt("id")
		if id <= 0 {
			c.Abort("404")
		}
		user := models.NewUser()
		userInfo, err := user.GetUserInfoById(id)
		if err != nil {
			c.Abort("404")
		}
		c.Data["id"] = id
		c.Data["userInfo"] = userInfo
		c.TplName = "admin/member-edit.html"
	}

}

func (c *AdminController) MemberDel() {
	id, _ := c.GetInt("id")
	if id < 0 {
		c.Json(300, "id有误")
	}
	user := models.NewUser()
	if err := user.UserDel(id); err != nil {

		c.Json(300, err.Error())
		return
	}
}

func (c *AdminController) Json(code int, msg string) {
	res := Response{
		Errcode: code,
		Message: msg,
	}
	c.Data["json"] = &res
	c.ServeJSON()
}

func (c *AdminController) BaseUrl() string {
	baseUrl := web.AppConfig.DefaultString("baseurl", "")
	if baseUrl != "" {
		if strings.HasSuffix(baseUrl, "/") {
			baseUrl = strings.TrimSuffix(baseUrl, "/")
		}
	} else {
		baseUrl = c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host
	}
	return baseUrl
}

func (c *AdminController) ArticleList() {
	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 10

	article := models.NewArticle()
	list, err := article.GetList(pageIndex, pageSize)
	if err != nil {

	}
	var articleList []*ArticleResData
	for _, v := range list {
		articleList = append(articleList, &ArticleResData{
			Id:         v.Id,
			Name:       v.Name,
			Content:    v.Content,
			CreateTime: time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05"),
		})
	}

	c.Data["data"] = articleList
	totalCount, _ := article.TotalCount()

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, int(totalCount), pageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	c.Data["totalCount"] = totalCount
	c.TplName = "admin/article-list.html"
}

func (c *AdminController) ArticleAdd() {

	valid := validation.Validation{}

	if c.Ctx.Input.IsPost() {
		username := c.GetString("name")

		content := c.GetString("content")

		if v := valid.Required(username, "username").Message("标题不能为空！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}
		if v := valid.Required(content, "password1").Message("内容不能为空！"); !v.Ok {
			c.Json(300, v.Error.Message)
			return
		}

		//添加数据库
		article := models.NewArticle()
		article.Name = username
		article.Content = content
		article.CreateTime = time.Now().Unix()
		article.UpdateTime = time.Now().Unix()
		if err := article.Add(); err != nil {
			c.Json(300, err.Error())
			return
		}
		c.Json(200, "添加成功")
	} else {

		c.TplName = "admin/article-add.html"
	}
}
