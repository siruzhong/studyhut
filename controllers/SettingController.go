package controllers

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"studyhut/constant"
	"studyhut/models"
	"studyhut/utils"
	"studyhut/utils/store"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type SettingController struct {
	BaseController
}

// Index 个人设置
func (this *SettingController) Index() {
	if this.Ctx.Input.IsPost() {
		email := strings.TrimSpace(this.GetString("email", ""))
		phone := strings.TrimSpace(this.GetString("phone"))
		wechatNO := strings.TrimSpace(this.GetString("wechat_no"))
		description := strings.TrimSpace(this.GetString("description"))
		nickname := strings.TrimSpace(this.GetString("nickname"))
		if email == "" {
			this.JsonResult(601, "邮箱不能为空")
		}

		if l := strings.Count(nickname, "") - 1; l < 2 || l > 20 {
			this.JsonResult(6004, "用户昵称限制在2-20个字符")
		}

		existMember := models.NewMember().FindByNickname(nickname, "member_id")
		if existMember.MemberId > 0 && this.Member.MemberId != existMember.MemberId {
			this.JsonResult(6004, "用户昵称已存在，请换一个")
		}

		member := this.Member
		member.Email = email
		member.Phone = phone
		member.WechatNO = wechatNO
		member.Description = description
		if err := member.Update(); err != nil {
			this.JsonResult(602, err.Error())
		}
		this.SetMember(*member)
		this.JsonResult(0, "ok")
	}
	this.Data["SeoTitle"] = "基本信息"
	this.Data["SettingBasic"] = true
	this.TplName = "setting/index.html"
}

// Password 修改密码
func (this *SettingController) Password() {
	if this.Ctx.Input.IsPost() {
		if this.Member.AuthMethod == constant.AuthMethodLDAP {
			this.JsonResult(6009, "当前用户不支持修改密码")
		}
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		password3 := this.GetString("password3")
		if password1 == "" {
			this.JsonResult(6003, "原密码不能为空")
		}

		if password2 == "" {
			this.JsonResult(6004, "新密码不能为空")
		}

		if count := strings.Count(password2, ""); count < 6 || count > 18 {
			this.JsonResult(6009, "密码必须在6-18字之间")
		}

		if password2 != password3 {
			this.JsonResult(6003, "确认密码不正确")
		}

		if ok, _ := utils.PasswordVerify(this.Member.Password, password1); !ok {
			this.JsonResult(6005, "原始密码不正确")
		}

		if password1 == password2 {
			this.JsonResult(6006, "新密码不能和原始密码相同")
		}

		pwd, err := utils.PasswordHash(password2)
		if err != nil {
			this.JsonResult(6007, "密码加密失败")
		}

		this.Member.Password = pwd
		if err := this.Member.Update(); err != nil {
			this.JsonResult(6008, err.Error())
		}

		this.JsonResult(0, "ok")
	}
	this.Data["SettingPwd"] = true
	this.Data["SeoTitle"] = "修改密码"
	this.TplName = "setting/password.html"
}

// Upload 上传图片
func (this *SettingController) Upload() {
	file, moreFile, err := this.GetFile("image-file")
	if err != nil {
		logs.Error("", err.Error())
		this.JsonResult(500, "读取文件异常")
	}
	defer file.Close()
	ext := filepath.Ext(moreFile.Filename) // 获取文件的拓展名
	// 文件拓展名只能为.png/.jpg/.gif/.jpeg
	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		this.JsonResult(500, "不支持的图片格式")
	}
	x1, _ := strconv.ParseFloat(this.GetString("x"), 10)
	y1, _ := strconv.ParseFloat(this.GetString("y"), 10)
	w1, _ := strconv.ParseFloat(this.GetString("width"), 10)
	h1, _ := strconv.ParseFloat(this.GetString("height"), 10)
	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)
	// 保存上传的图片
	fileName := strconv.FormatInt(time.Now().UnixNano(), 16)
	filePath := filepath.Join("uploads", time.Now().Format("2006/01"), fileName+ext)
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	err = this.SaveToFile("image-file", filePath)
	if err != nil {
		logs.Error("", err)
		this.JsonResult(500, "图片保存失败")
	}
	// 剪切图片
	subImg, err := utils.ImageCopyFromFile(filePath, x, y, width, height)
	if err != nil {
		logs.Error("ImageCopyFromFile => ", err)
		this.JsonResult(6001, "头像剪切失败")
	}
	// 删除原图片
	os.Remove(filePath)
	// 保存剪切后的图片
	filePath = filepath.Join("uploads", time.Now().Format("2006/01"), fileName+ext)
	utils.ImageResize(subImg, 120, 120)
	err = utils.SaveImage(filePath, subImg)
	if err != nil {
		logs.Error("保存文件失败 => ", err.Error())
		this.JsonResult(500, "保存文件失败")
	}
	url := "/" + strings.Replace(filePath, "\\", "/", -1)
	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}
	switch utils.StoreType {
	case constant.StoreOss: //oss存储
		if err := store.ModelStoreOss.MoveToOss("."+url, strings.TrimLeft(url, "/"), true, false); err != nil {
			beego.Error(err.Error())
		} else {
			url = strings.TrimRight(beego.AppConfig.String("oss::Domain"), "/ ") + url
		}
	case constant.StoreCos: //cos存储
		if err := store.ModelStoreCos.MoveToCos("."+url, strings.TrimLeft(url, "/"), true, false); err != nil {
			beego.Error(err.Error())
		} else {
			url = strings.TrimRight(beego.AppConfig.String("cos::Domain"), "/ ") + url
		}
	case constant.StoreLocal: //本地存储
		if err := store.ModelStoreLocal.MoveToStore("."+url, strings.TrimLeft(url, "./")); err != nil {
			beego.Error(err.Error())
		} else {
			url = "/" + strings.TrimLeft(url, "./")
		}
	}
	if member, err := models.NewMember().Find(this.Member.MemberId); err == nil {
		avatar := member.Avatar
		member.Avatar = url
		err = member.Update("avatar")
		if err != nil {
			this.JsonResult(60001, "更新头像失败")
		}
		avatar = strings.TrimLeft(avatar, "./")
		if strings.HasPrefix(avatar, beego.AppConfig.String("cos::Domain")) {
			store.ModelStoreCos.DelFromCos(strings.Replace(avatar, beego.AppConfig.String("cos::Domain"), "", 1)) // cos上删除原头像
		} else if strings.HasPrefix(avatar, beego.AppConfig.String("oos::Domain")) {
			store.ModelStoreCos.DelFromCos(strings.Replace(avatar, beego.AppConfig.String("oss::Domain"), "", 1)) // oss上删除原头像
		} else {
			os.Remove(avatar) // 本地删除原头像
		}
		this.SetMember(*member)
	}
	this.JsonResult(0, "ok", url)
}
