package controllers

import (
	"github.com/astaxie/beego"
	"programming-learning-platform/constant"
	"programming-learning-platform/models"
	"programming-learning-platform/utils"
)

// StarController 收藏控制器
type StarController struct {
	BaseController
}

// List 我的收藏列表
func (this *StarController) List() {
	page, _ := this.GetInt("page")
	cid, _ := this.GetInt("cid")
	if page < 1 {
		page = 1
	}
	sort := this.GetString("sort", "read")
	cnt, books, _ := new(models.Star).List(this.Member.MemberId, page, constant.PageSize, cid, sort)
	if cnt > 1 {
		this.Data["PageHtml"] = utils.NewPaginations(constant.RollPage, int(cnt), constant.PageSize, page, beego.URLFor("StarController.List"), "")
	}
	this.Data["Pid"] = 0
	cates := models.NewCategory().CategoryOfUserCollection(this.Member.MemberId)
	for _, cate := range cates {
		if cate.Id == cid {
			if cate.Pid == 0 {
				this.Data["Pid"] = cate.Id
			} else {
				this.Data["Pid"] = cate.Pid
			}
		}
	}
	this.Data["Books"] = books
	this.Data["Sort"] = sort
	this.Data["SettingStar"] = true
	this.Data["SeoTitle"] = "我的收藏 - " + this.Sitename
	this.TplName = "setting/star.html"
	this.Data["Cid"] = cid
	this.Data["Cates"] = cates
}
