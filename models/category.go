package models

import (
	"errors"
	"studyhut/constant"
	"studyhut/utils/store"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"studyhut/utils"
)

var tableCategory = "category"

// Category 分类
type Category struct {
	Id     int    `json:"id"`                                    // 自增主键
	Pid    int    `json:"pid"`                                   // 父分类id
	Title  string `orm:"size(30);unique" json:"title,omitempty"` // 分类名称
	Intro  string `json:"intro,omitempty"`                       // 介绍
	Icon   string `json:"icon,omitempty"`                        // 分类icon
	Cnt    int    `json:"cnt,omitempty"`                         // 分类下的书籍统计
	Sort   int    `json:"sort,omitempty"`                        // 排序
	Status bool   `json:"status,omitempty"`                      // 分类状态，true表示显示，否则表示隐藏
}

// NewCategory 创建分类
func NewCategory() *Category {
	return &Category{}
}

// AddCategory 新增分类
func (this *Category) AddCategory(pid int, cates string) (err error) {
	slice := strings.Split(cates, "\n")
	if len(slice) == 0 {
		return
	}
	o := orm.NewOrm()
	for _, item := range slice {
		if item = strings.TrimSpace(item); item != "" {
			var cate = Category{
				Pid:    pid,
				Title:  item,
				Status: true,
			}
			if cnt, _ := o.QueryTable(this).Filter("title", cate.Title).Filter("pid", cate.Pid).Count(); cnt == 0 {
				_, err = orm.NewOrm().Insert(&cate)
			}
		}
	}
	return
}

// DelCategory 删除分类（如果分类下的书籍不为0，则不允许删除）
func (this *Category) DelCategory(id int) (err error) {
	var cate = Category{Id: id}

	o := orm.NewOrm()
	if err = o.Read(&cate); cate.Cnt > 0 { //当前分类下书籍数量不为0，不允许删除
		return errors.New("删除失败，当前分类下的问下书籍不为0，不允许删除")
	}

	if _, err = o.Delete(&cate, "id"); err != nil {
		return
	}
	_, err = o.QueryTable(tableCategory).Filter("pid", id).Delete()
	if err != nil { //删除分类图标
		return
	}

	switch utils.StoreType {
	case constant.StoreOss:
		store.ModelStoreOss.DelFromOss(cate.Icon)
	case constant.StoreCos:
		store.ModelStoreCos.DelFromCos(cate.Icon)
	case constant.StoreLocal:
		store.ModelStoreLocal.DelFiles(cate.Icon)
	}
	return
}

// GetAllCategory 查询所有分类
// @param	pid		-1表示不限(即查询全部)，否则表示查询指定pid的分类
// @param	status	-1表示不限状态(即查询所有状态的分类)，0表示关闭状态，1表示启用状态
func (this *Category) GetAllCategory(pid int, status int) (cates []Category, err error) {
	qs := orm.NewOrm().QueryTable(tableCategory)
	if pid > -1 {
		qs = qs.Filter("pid", pid)
	}
	if status == 0 || status == 1 {
		qs = qs.Filter("status", status)
	}
	_, err = qs.OrderBy("-status", "sort", "title").All(&cates)
	return
}

// UpdateByField 根据字段更新内容
func (this *Category) UpdateByField(id int, field, val string) (err error) {
	_, err = orm.NewOrm().QueryTable(tableCategory).Filter("id", id).Update(orm.Params{field: val})
	return
}

// Find 查询单个分类
func (this *Category) Find(id int) (cate Category) {
	cate.Id = id
	orm.NewOrm().Read(&cate)
	return cate
}

// CategoryOfUserCollection 用户收藏了的书籍的分类
func (this *Category) CategoryOfUserCollection(uid int, forAPI ...bool) (cates []Category) {
	order := " ORDER BY c.sort asc,c.title asc "
	if len(forAPI) > 0 && forAPI[0] {
		order = " ORDER BY c.title asc "
	}
	sql := `
		SELECT c.id,c.pid,c.title
		FROM book_category bc
			LEFT JOIN star s ON s.bid = bc.book_id
			LEFT JOIN category c ON c.id = bc.category_id
		WHERE s.uid = ? 
		GROUP BY c.id ` + order
	if _, err := orm.NewOrm().Raw(sql, uid).QueryRows(&cates); err != nil {
		beego.Error(err.Error())
	}
	return
}
