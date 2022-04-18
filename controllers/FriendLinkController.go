package controllers

import (
	"strings"
	"studyhut/models"
	"studyhut/utils"
)

// FriendLinkController 友情链接控制器
type FriendLinkController struct {
	BaseController
}

// List 展示所有友链
func (this *FriendLinkController) List() {
	friendlinks := new(models.FriendLink).GetList(false)
	for idx, friendlink := range friendlinks {
		if strings.TrimSpace(friendlink.Pic) == "" { // 赋值为默认图片
			friendlink.Pic = "/static/images/icon.png"
		} else {
			friendlink.Pic = utils.ShowImg(friendlink.Pic)
		}
		friendlinks[idx] = friendlink
	}
	this.Data["Friendlinks"] = friendlinks
	this.TplName = "friendlink/list.html"
	this.GetSeoByPage("friendlink_list", map[string]string{
		"title":       "友链列表",
		"keywords":    "友链列表",
		"description": "友情链接列表，优质的网站推荐",
	})
}
