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
		"title":       "书籍分类",
		"keywords":    "文档托管,在线创作,文档在线管理,在线知识管理,文档托管平台,在线写书,文档在线转换,在线编辑,在线阅读,开发手册,api手册,文档在线学习,技术文档,在线编辑",
		"description": this.Sitename + "专注于文档在线写作、协作、分享、阅读与托管，让每个人更方便地发布、分享和获得知识。",
	})
	this.Data["IsCate"] = true
	this.Data["Friendlinks"] = new(models.FriendLink).GetList(false)
	this.Data["Recommends"], _, _ = models.NewBook().HomeData(1, 12, constant.OrderLatestRecommend, "", 0)
	this.Data["SHOW_CATEGORY_INDEX"] = "true"
	this.Data["Cates"], _ = models.NewCategory().GetAllCategory(-1, -1)
	this.TplName = "cates/list.html"
}
