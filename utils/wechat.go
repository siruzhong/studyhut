package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type WechatToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

var wechatToken = &WechatToken{}

// GetAccessToken 获取access_token
func GetAccessToken() (token *WechatToken) {
	now := time.Now().Unix()
	if now < wechatToken.ExpiresIn-600 { //提前10分钟失效
		return wechatToken
	}
	appId := beego.AppConfig.String("appId")
	appSecret := beego.AppConfig.String("appSecret")
	api := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", appId, appSecret)
	req := httplib.Get(api).SetTimeout(10*time.Second, 10*time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := req.String()
	if err != nil {
		beego.Error(err.Error())
	}
	token = &WechatToken{}
	json.Unmarshal([]byte(resp), token)
	token.ExpiresIn = now + token.ExpiresIn
	wechatToken = token
	return
}
