package controllers

import (
	"github.com/astaxie/beego"
	"studyhut/constant"
	"studyhut/models"
)

// CateController 分类控制器
type CateController struct {
	BaseController
}

// Index 分类首页
func (this *CateController) Index() {
	cid, _ := this.GetInt("cid")
	if cid > 0 {
		this.Redirect(beego.URLFor("HomeController.Index")+this.Ctx.Request.RequestURI, 302)
	}
	this.List()
}

// List 分类展示
func (this *CateController) List() {
	if cates, err := new(models.Category).GetAllCategory(-1, 1); err == nil {
		this.Data["Cates"] = cates
	} else {
		beego.Error(err.Error())
	}
	this.GetSeoByPage("cate", map[string]string{
		"title":       "首页",
		"keywords":    "IT技术、资源整合、在线学习、交流分享、内容创作",
		"description": this.Sitename + "一个在线IT技术资源整合、在线学习、交流分享的站点。每一名用户都是内容的创造者，分享你认为优质的资源，让我们一起学习！一起进步！",
	})
	this.Data["IsCate"] = true
	this.Data["Friendlinks"] = new(models.FriendLink).GetList(false)
	this.Data["Recommends"], _, _ = models.NewBook().HomeData(1, 12, constant.OrderLatestRecommend, "", 0)
	this.Data["SHOW_CATEGORY_INDEX"] = "true"
	this.Data["Cates"], _ = models.NewCategory().GetAllCategory(-1, -1)
	this.TplName = "cates/list.html"
}
