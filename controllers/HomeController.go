package controllers

import (
	"math"
	"strconv"
	"strings"
	"studyhut/constant"

	"github.com/astaxie/beego"
	"studyhut/models"
	"studyhut/utils"
)

type HomeController struct {
	BaseController
}

func (this *HomeController) Index() {
	var (
		tab       string
		cid       int // 分类，如果只是一级分类，则忽略；二级分类，则根据二级分类查找内容
		urlPrefix = beego.URLFor("HomeController.Index")
		cate      models.Category
		lang      = this.GetString("lang")
		tabName   = map[string]string{"recommend": "网站推荐", "latest": "最新发布", "popular": "热门书籍"}
	)

	tab = strings.ToLower(this.GetString("tab"))
	switch tab {
	case "recommend", "popular", "latest":
	default:
		tab = "latest"
	}

	ModelCate := new(models.Category)
	cates, _ := ModelCate.GetAllCategory(-1, 1)
	cid, _ = this.GetInt("cid")
	pid := cid
	if cid > 0 {
		for _, item := range cates {
			if item.Id == cid {
				if item.Pid > 0 {
					pid = item.Pid
				}
				this.Data["Cate"] = item
				cate = item
				break
			}
		}
	}
	this.Data["Cates"] = cates
	this.Data["Cid"] = cid
	this.Data["Pid"] = pid
	this.TplName = "home/index.html"
	this.Data["IsHome"] = true

	pageIndex, _ := this.GetInt("page", 1)
	// 每页显示24个，为了兼容Pad、mobile、PC
	pageSize := 24
	books, totalCount, err := models.NewBook().HomeData(pageIndex, pageSize, models.BookOrder(tab), lang, cid)
	if err != nil {
		beego.Error(err)
		this.Abort("404")
	}
	if totalCount > 0 {
		urlSuffix := "&tab=" + tab
		if cid > 0 {
			urlSuffix = urlSuffix + "&cid=" + strconv.Itoa(cid)
		}
		urlSuffix = urlSuffix + "&lang=" + lang
		html := utils.NewPaginations(constant.RollPage, totalCount, pageSize, pageIndex, urlPrefix, urlSuffix)
		this.Data["PageHtml"] = html
	} else {
		this.Data["PageHtml"] = ""
	}

	this.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	this.Data["Lists"] = books
	this.Data["Tab"] = tab
	this.Data["Lang"] = lang
	title := this.Sitename

	desc := this.Sitename + "一个在线IT技术资源整合、在线学习、交流分享的站点。每一名用户都是内容的创造者，分享你认为优质的资源，让我们一起学习！一起进步！"
	if cid > 0 {
		title = "[发现页] " + cate.Title + " - " + tabName[tab] + " - " + title
		if strings.TrimSpace(cate.Intro) != "" {
			desc = cate.Title + "，" + cate.Intro + " - " + this.Sitename
		}
	} else {
		title = "发现"
	}

	this.Data["Cate"] = cate

	this.GetSeoByPage("index", map[string]string{
		"title":       title,
		"keywords":    "IT技术、资源整合、在线学习、交流分享、内容创作",
		"description": desc,
	})
}
