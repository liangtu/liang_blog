package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"liang_blog/models"
	pagination "liang_blog/utiils"
	"math"
	"strings"
	"time"
)

type IndexController struct {
	web.Controller
}

func (c *IndexController) URLMapping() {
	c.Mapping("Index", c.Index)
	c.Mapping("Info", c.Info)
}

func (c *IndexController) Index() {

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
	c.TplName = "index/index.html"
}

func (c *IndexController) Info() {
	id, _ := c.GetInt("id")
	if id <= 0 {
		c.Abort("404")
	}
	a := models.NewArticle()
	article, err := a.GetArticleInfoById(id)
	if err != nil {
		c.Abort("404")
	}
	c.Data["article"] = article
	c.TplName = "index/info.html"
}

func (c *IndexController) BaseUrl() string {
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
