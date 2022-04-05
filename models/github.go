package models

import (
	"github.com/astaxie/beego/orm"
	"studyhut/oauth"
)

var ModelGithub = new(Github)

type Github struct {
	oauth.GithubUser
}

// GetUserByGithubId 根据GithubId获取用户的GitHub数据。这里可以查询用户是否绑定了或者数据是否在库中存在
func (this *Github) GetUserByGithubId(id int, cols ...string) (user Github, err error) {
	qs := orm.NewOrm().QueryTable("github").Filter("id", id)
	if len(cols) > 0 {
		err = qs.One(&user, cols...)
	} else {
		err = qs.One(&user)
	}
	return
}

// Bind 绑定用户
func (this *Github) Bind(githubId, memberId interface{}) (err error) {
	_, err = orm.NewOrm().QueryTable("github").Filter("id", githubId).Update(orm.Params{"member_id": memberId})
	return
}
