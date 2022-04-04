package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"programming-learning-platform/constant"
	"programming-learning-platform/utils/store"
	"strings"

	"github.com/astaxie/beego"
	"programming-learning-platform/models"
	"programming-learning-platform/utils"
)

// StaticController 静态文件控制器
type StaticController struct {
	beego.Controller
	OssDomain string
}

func (this *StaticController) Prepare() {
	this.OssDomain = strings.TrimRight(beego.AppConfig.String("oss::Domain"), "/ ")
}

func (this *StaticController) APP() {
	link := strings.TrimSpace(models.GetOptionValue("APP_PAGE", ""))
	if link != "" {
		this.Redirect(link, 302)
	}
	this.Abort("404")
}

// Uploads 查询上传的静态资源
// 如果是音频和视频文件，需要根据后台设置而判断是否加密处理
// 如果使用了OSS存储，则需要将文件处理好
func (this *StaticController) Uploads() {
	file := strings.TrimLeft(this.GetString(":splat"), "./")
	path := strings.ReplaceAll(filepath.Join("uploads", file), "\\", "/")
	if this.isMedia(path) { // 签名验证
		sign := this.GetString("sign")
		if !this.isValidSign(sign, path) {
			// 签名验证不通过，需要再次验证书籍是否是用户的（针对编辑状态）
			if !this.isBookOwner() {
				this.Abort("404")
				return
			}
		}
	}
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, path)
}

// StaticFile 静态文件，这个加在路由的最后
func (this *StaticController) StaticFile() {
	file := this.GetString(":splat")
	if strings.HasPrefix(file, ".well-known") || file == "sitemap.xml" {
		http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, file)
		return
	}

	file = strings.ReplaceAll(strings.TrimLeft(file, "./"), "\\", "/")
	//path := filepath.Join(utils.VirtualRoot, file)
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, file)
}

// ProjectsFile 书籍静态文件
func (this *StaticController) ProjectsFile() {
	if utils.StoreType != constant.StoreOss {
		this.Abort("404")
	}

	object := filepath.Join("projects/", strings.TrimLeft(this.GetString(":splat"), "./"))
	object = strings.ReplaceAll(object, "\\", "/")

	// 不是音频和视频，直接跳转
	if !this.isMedia(object) {
		this.Redirect(this.OssDomain+"/"+object, 302)
		return
	}

	// 签名验证
	sign := this.GetString("sign")
	if !this.isValidSign(sign, object) {
		// 签名验证不通过，需要再次验证书籍是否是用户的（针对编辑状态）
		if !this.isBookOwner() {
			this.Abort("404")
			return
		}
	}

	if bucket, err := store.ModelStoreOss.GetBucket(); err == nil {
		object, _ = bucket.SignURL(object, http.MethodGet, utils.MediaDuration)
		if slice := strings.Split(object, "/"); len(slice) > 2 {
			object = strings.Join(slice[3:], "/")
		}
	}
	this.Redirect(this.OssDomain+"/"+object, 302)
}

// isMedia 是否是音视频
func (this *StaticController) isMedia(path string) (yes bool) {
	var videoOK, audioOK bool
	ext := strings.ToLower(filepath.Ext(path))
	_, videoOK = utils.VideoExt.Load(ext) // 判断是否为视频拓展名
	_, audioOK = utils.AudioExt.Load(ext) // 判断是否为音频拓展名
	return audioOK || videoOK
}

// isBookOwner 是否是书籍项目所有人(书籍项目所有人，可以直链播放音视频)
func (this *StaticController) isBookOwner() (yes bool) {
	memberID := 0
	// 从session中获取用户信息
	if member, ok := this.GetSession(constant.LoginSessionName).(models.Member); ok {
		memberID = member.MemberId
	}

	if memberID <= 0 {
		// 如果Cookie中存在登录信息，从cookie中获取用户信息
		if cookie, ok := this.GetSecureCookie(utils.GetAppKey(), "login"); ok {
			var remember CookieRemember
			if err := utils.Decode(cookie, &remember); err == nil {
				memberID = remember.MemberId
			}
		}
	}
	if memberID <= 0 {
		return
	}

	referer := this.Ctx.Request.Referer()
	if referer == "" {
		return
	}

	bookIdentify := ""
	if u, err := url.Parse(referer); err == nil {
		fmt.Println(u.Path)
		if slice := strings.Split(u.Path, "/"); len(slice) >= 3 && slice[1] == "api" {
			bookIdentify = slice[2]
		}
	}

	if bookIdentify == "" {
		return
	}

	bookID := 0
	if book, err := models.NewBook().FindByIdentify(bookIdentify, "book_id"); err == nil {
		bookID = book.BookId
	}
	if bookID <= 0 {
		return
	}

	if r, err := models.NewRelationship().FindByBookIdAndMemberId(bookID, memberID); err == nil && r.RelationshipId > 0 {
		return true
	}

	return false
}

// isValidSign 是否是合法的签名（针对音频和视频，签名不可用的时候再验证用户有没有登录，用户登录了再验证用户是不是书籍所有人）
func (this *StaticController) isValidSign(sign, path string) bool {
	signPath, err := utils.ParseMediaSign(sign)
	if err != nil {
		return false
	}
	return signPath == path
}
