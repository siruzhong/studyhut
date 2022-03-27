package models

import (
	"fmt"

	"github.com/TruthHun/gotil/cryptil"
	"github.com/astaxie/beego/orm"
)

// Wechat 第三方微信接口
type Wechat struct {
	Id        int
	MemberId  int    //绑定的用户id
	Openid    string `orm:"unique;size(50)"`
	Unionid   string `orm:"size(50)"`
	AvatarURL string `orm:"column(avatar_url)"`
	Nickname  string `orm:"size(30)"`
	SessKey   string `orm:"size(50);unique"`
}

// NewWechat 创建微信实体
func NewWechat() *Wechat {
	return &Wechat{}
}

// GetUserByOpenid 根据openid获取用户的微信数据
func (this *Wechat) GetUserByOpenid(openid string, cols ...string) (user Wechat, err error) {
	// 查询用户的微信数据是否在数据库中存在
	qs := orm.NewOrm().QueryTable(this).Filter("openid", openid)
	if len(cols) > 0 {
		err = qs.One(&user, cols...)
	} else {
		err = qs.One(&user)
	}
	return
}

// GetUserBySess 根据SessKey获取用户的微信数据
func (this *Wechat) GetUserBySess(sessKey string, cols ...string) (user Wechat, err error) {
	qs := orm.NewOrm().QueryTable(this).Filter("sess_key", sessKey)
	if len(cols) > 0 {
		err = qs.One(&user, cols...)
	} else {
		err = qs.One(&user)
	}
	return
}

// Insert 插入用户微信数据到数据库
func (this *Wechat) Insert() (err error) {
	o := orm.NewOrm()
	exist := &Wechat{}
	o.QueryTable(this).Filter("openid", this.Openid).One(exist)
	if exist.Id > 0 {
		exist.SessKey = this.SessKey
		_, err = o.Update(exist)
	} else {
		_, err = o.Insert(this)
	}
	return
}

// Bind 绑定用户与微信
func (this *Wechat) Bind(openid, memberId interface{}) (err error) {
	_, err = orm.NewOrm().QueryTable(this).Filter("openid", openid).Filter("member_id", 0).Update(orm.Params{"member_id": memberId, "sess_key": cryptil.Md5Crypt(fmt.Sprint(openid))})
	return
}
