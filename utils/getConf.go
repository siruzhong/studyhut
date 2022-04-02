package utils

import (
	"github.com/astaxie/beego"
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
