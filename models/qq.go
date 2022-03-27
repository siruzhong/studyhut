package models

import "github.com/astaxie/beego/orm"

// QQ 第三方qq接口
type QQ struct {
	Id        int
	OpenId    string
	MemberId  int
	Name      string `orm:"size(50)"` //昵称
	Gender    string `orm:"size(5)"`
	AvatarURL string `orm:"column(avatar_url)"` //用户头像链接
}

// ModelQQ qq实体
var ModelQQ = new(QQ)

// TableName 获取表名
func (this *QQ) TableName() string {
	return "qq"
}

// GetUserByOpenid 根据openid获取用户的qq数据
func (this *QQ) GetUserByOpenid(openid string, cols ...string) (user QQ, err error) {
	// 查询用户的qq数据是否在数据库中存在
	qs := orm.NewOrm().QueryTable("qq").Filter("openid", openid)
	if len(cols) > 0 {
		err = qs.One(&user, cols...)
	} else {
		err = qs.One(&user)
	}
	return
}

// Bind 绑定用户与qq
func (this *QQ) Bind(openid, memberId interface{}) (err error) {
	_, err = orm.NewOrm().QueryTable("qq").Filter("openid", openid).Update(orm.Params{"member_id": memberId})
	return
}
