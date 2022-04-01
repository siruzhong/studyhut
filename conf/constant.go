package conf

import (
	"github.com/astaxie/beego"
)

// 用户角色
const (
	MemberSuperRole   = 0 // 超级管理员
	MemberAdminRole   = 1 // 普通管理员
	MemberGeneralRole = 2 // 读者
	MemberEditorRole  = 3 // 作者（可以创建书籍）
)

// 书籍角色
const (
	BookFounder  = 0 // 创始人
	BookAdmin    = 1 // 管理者
	BookEditor   = 2 // 编辑者
	BookObserver = 3 // 观察者
)

// 广告位置
const (
	AdsPositionBeforeFriendLink        = "global-before-friend-link"
	AdsPositionGlobalFooter            = "global-footer"
	AdsPositionUnderLatestRecommend    = "index-under-latest-recommend"
	AdsPositionSearchRight             = "search-right"
	AdsPositionSearchTop               = "search-top"
	AdsPositionSearchBottom            = "search-bottom"
	AdsPositionUnderBookName           = "intro-under-book-name"
	AdsPositionBeforeMenu              = "intro-before-menu"
	AdsPositionBeforeRelatedBooks      = "intro-before-related-books"
	AdsPositionUnderExploreNav         = "explore-under-nav"
	AdsPositionBeforeExplorePagination = "explore-before-pagination"
	AdsPositionUnderExplorePagination  = "explore-under-pagination"
	AdsPositionContentTop              = "content-top"
	AdsPositionContentBottom           = "content-bottom"
)

const (
	LoginSessionName = "LoginSessionName" // 登录用户的Session名

	RegexpEmail   = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	RegexpAccount = `^[a-zA-z0-9_]{2,50}$` // 允许用户名中出现点号

	PageSize = 10 // 默认分页条数
	RollPage = 4  // 展示分页的个数

	AuthMethodLocal = "local" // 本地账户校验
	AuthMethodLDAP  = "ldap"  // LDAP用户校验
)

// 字符串类型
const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)


// GetAppKey 获取app_key
func GetAppKey() string {
	return beego.AppConfig.DefaultString("app_key", "godoc")
}

// GetDatabasePrefix 获取数据库前缀
func GetDatabasePrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "")
}

// GetDefaultAvatar 获取默认头像
func GetDefaultAvatar() string {
	return beego.AppConfig.DefaultString("avatar", "/static/images/headimgurl.jpg")
}

// GetTokenSize 获取阅读令牌长度
func GetTokenSize() int {
	return beego.AppConfig.DefaultInt("token_size", 12)
}

// GetDefaultCover 获取默认文档封面
func GetDefaultCover() string {
	return beego.AppConfig.DefaultString("cover", "/static/images/book.jpg")
}
