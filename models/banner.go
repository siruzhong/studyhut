package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"programming-learning-platform/utils"
)

// Banner 横幅
type Banner struct {
	Id        int       `json:"id"`                        // 横幅id
	Type      string    `orm:"size(30);index" json:"type"` // 横幅类型
	Title     string    `json:"title" orm:"size(100)"`     // 横幅名称
	Link      string    `json:"link"`                      // 横幅链接
	Image     string    `json:"image"`                     // 横幅图片
	Sort      int       `json:"sort"`                      // 排序号
	Status    bool      `json:"status"`                    // 状态(1为启用,0为未启用)
	CreatedAt time.Time `json:"created_at"`                // 创建时间
}

// NewBanner 创建横幅
func NewBanner() *Banner {
	return &Banner{}
}

// Lists 展示指定类型的横幅(根据sort和id降序排列)
func (m *Banner) Lists(t string) (banners []Banner, err error) {
	_, err = orm.NewOrm().QueryTable(m).Filter("type", t).Filter("status", true).OrderBy("-sort", "-id").All(&banners)
	if err == orm.ErrNoRows {
		err = nil
	}
	return
}

// All 查询所有横幅(根据sort和status降序排列)
func (m *Banner) All() (banners []Banner, err error) {
	_, err = orm.NewOrm().QueryTable(m).OrderBy("-sort", "-status").All(&banners)
	if err == orm.ErrNoRows {
		err = nil
	}
	return
}

// Update 更新横幅
func (m *Banner) Update(id int, field string, value interface{}) (err error) {
	_, err = orm.NewOrm().QueryTable(m).Filter("id", id).Update(orm.Params{field: value})
	if err == orm.ErrNoRows {
		err = nil
	}
	return
}

// Delete 删除横幅
func (m *Banner) Delete(id int) (err error) {
	var banner Banner
	q := orm.NewOrm().QueryTable(m).Filter("id", id)
	q.One(&banner)
	if banner.Id > 0 {
		_, err = q.Delete()
		if err == nil {
			utils.DeleteFile(banner.Image)
		}
	}
	return
}
