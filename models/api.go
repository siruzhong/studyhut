package models

import (
	"strings"
	"studyhut/constant"

	"github.com/astaxie/beego"
	"studyhut/utils"
)

var staticDomain string

func initAPI() {
	if strings.ToLower(utils.StoreType) == constant.StoreOss {
		staticDomain = strings.TrimSpace(beego.AppConfig.String("oss::Domain"))
	}

	if strings.ToLower(utils.StoreType) == constant.StoreCos {
		staticDomain = strings.TrimSpace(beego.AppConfig.String("cos::Domain"))
	}

	if strings.TrimRight(staticDomain, "/") == "" {
		staticDomain = beego.AppConfig.DefaultString("static_domain", "")
	}

	staticDomain = strings.TrimRight(staticDomain, "/") + "/"
}
