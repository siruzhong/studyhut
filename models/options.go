package models

import (
	"strings"
	"studyhut/utils"
	"sync"

	"github.com/astaxie/beego/orm"
)

// Option 配置项
type Option struct {
	OptionId    int    `orm:"column(option_id);pk;auto;unique;" json:"option_id"`       // 配置id
	OptionTitle string `orm:"column(option_title);size(500)" json:"option_title"`       // 配置标题
	OptionName  string `orm:"column(option_name);unique;size(80)" json:"option_name"`   // 配置名称
	OptionValue string `orm:"column(option_value);type(text);null" json:"option_value"` // 配置值
	Remark      string `orm:"column(remark);type(text);null" json:"remark"`             // 描述
}

// NewOption 创建配置项
func NewOption() *Option {
	return &Option{}
}

// TableName 获取数据表名
func (m *Option) TableName() string {
	return "options"
}

// TableNameWithPrefix 获取带前缀对数据表名
func (m *Option) TableNameWithPrefix() string {
	return utils.GetDatabasePrefix() + m.TableName()
}

// optionCache 配置项缓存
var optionCache sync.Map

// initOptionCache 初始化配置缓存
func initOptionCache() {
	opts, _ := NewOption().All()
	for _, opt := range opts {
		optionCache.Store(opt.OptionName, opt)
		optionCache.Store(opt.OptionId, opt)
	}
}

// Init 初始化配置项
func (m *Option) Init() error {
	o := orm.NewOrm()
	options := []Option{
		{
			OptionValue: "true",
			OptionName:  "ENABLED_REGISTER",
			OptionTitle: "是否启用注册",
		},
		{
			OptionValue: "100",
			OptionName:  "ENABLE_DOCUMENT_HISTORY",
			OptionTitle: "版本控制",
		}, {
			OptionValue: "true",
			OptionName:  "ENABLED_CAPTCHA",
			OptionTitle: "是否启用验证码",
		}, {
			OptionValue: "true",
			OptionName:  "ENABLE_ANONYMOUS",
			OptionTitle: "启用匿名访问",
		}, {
			OptionValue: "BookStack",
			OptionName:  "SITE_NAME",
			OptionTitle: "站点名称",
		}, {
			OptionValue: "",
			OptionName:  "ICP",
			OptionTitle: "网站备案",
		}, {
			OptionValue: "",
			OptionName:  "TONGJI",
			OptionTitle: "站点统计",
		}, {
			OptionValue: "true",
			OptionName:  "SPIDER",
			OptionTitle: "采集器，是否只对管理员开放",
		}, {
			OptionValue: "false",
			OptionName:  "SHOW_CATEGORY_INDEX",
			OptionTitle: "首页是否显示分类索引",
		}, {
			OptionValue: "false",
			OptionName:  "ELASTICSEARCH_ON",
			OptionTitle: "是否开启全文搜索",
		}, {
			OptionValue: "http://localhost:9200/",
			OptionName:  "ELASTICSEARCH_HOST",
			OptionTitle: "ElasticSearch Host",
		}, {
			OptionValue: "book",
			OptionName:  "DEFAULT_SEARCH",
			OptionTitle: "默认搜索",
		}, {
			OptionValue: "50",
			OptionName:  "SEARCH_ACCURACY",
			OptionTitle: "搜索精度",
		}, {
			OptionValue: "true",
			OptionName:  "LOGIN_QQ",
			OptionTitle: "是否允许使用QQ登录",
		}, {
			OptionValue: "true",
			OptionName:  "LOGIN_GITHUB",
			OptionTitle: "是否允许使用Github登录",
		}, {
			OptionValue: "true",
			OptionName:  "LOGIN_GITEE",
			OptionTitle: "是否允许使用码云登录",
		}, {
			OptionValue: "1",
			OptionName:  "RELATE_BOOK",
			OptionTitle: "是否开始关联书籍",
		}, {
			OptionValue: "true",
			OptionName:  "ALL_CAN_WRITE_BOOK",
			OptionTitle: "是否都可以创建书籍",
		}, {
			OptionValue: "true",
			OptionName:  "CLOSE_OPEN_SOURCE_LINK",
			OptionTitle: "是否关闭开源书籍入口",
		}, {
			OptionValue: "X-Real-Ip",
			OptionName:  "REAL_IP_FIELD",
			OptionTitle: "request中获取访客真实IP的header",
		}, {
			OptionValue: "",
			OptionName:  "APP_PAGE",
			OptionTitle: "手机APP下载单页",
		}, {
			OptionValue: "false",
			OptionName:  "HIDE_TAG",
			OptionTitle: "是否隐藏标签在导航栏显示",
		}, {
			OptionValue: "",
			OptionName:  "DOWNLOAD_LIMIT",
			OptionTitle: "是否需要登录才能下载电子书",
		}, {
			OptionValue: "false",
			OptionName:  "AUTO_HTTPS",
			OptionTitle: "图片链接HTTP转HTTPS",
		}, {
			OptionValue: "5",
			OptionName:  "SIGN_BASIC_REWARD",
			OptionTitle: "用户每次签到基础奖励阅读时长(秒)",
		}, {
			OptionValue: "0",
			OptionName:  "SIGN_CONTINUOUS_REWARD",
			OptionTitle: "用户连续签到奖励阅读时长(秒)",
		}, {
			OptionValue: "0",
			OptionName:  "SIGN_CONTINUOUS_MAX_REWARD",
			OptionTitle: "连续签到奖励阅读时长上限(秒)",
		},
		{
			OptionValue: "0",
			OptionName:  "READING_MIN_INTERVAL",
			OptionTitle: "内容最小阅读计时间隔(秒)",
		},
		{
			OptionValue: "600",
			OptionName:  "READING_MAX_INTERVAL",
			OptionTitle: "内容最大阅读计时间隔(秒)",
		},
		{
			OptionValue: "1200",
			OptionName:  "READING_INVALID_INTERVAL",
			OptionTitle: "内容阅读无效计时间隔(秒)",
		},
		{
			OptionValue: "600",
			OptionName:  "READING_INTERVAL_MAX_REWARD",
			OptionTitle: "内容阅读计时间隔最大奖励(秒)",
		},
		{
			OptionValue: "false",
			OptionName:  "COLLAPSE_HIDE",
			OptionTitle: "目录是否默认收起",
		},
		{
			OptionValue: "",
			OptionName:  "FORBIDDEN_REFERER",
			OptionTitle: "禁止的Referer",
		}, {
			OptionValue: "1",
			OptionName:  "DOWNLOAD_INTERVAL",
			OptionTitle: "每阅读多少秒可以下载一个电子书",
		},
	}
	for _, op := range options {
		// 不存在则插入
		if !o.QueryTable(m.TableNameWithPrefix()).Filter("option_name", op.OptionName).Exist() {
			if _, err := o.Insert(&op); err != nil {
				return err
			}
		}
	}
	initOptionCache()
	return nil
}

// Find 根据id查找配置项
func (p *Option) Find(id int) (*Option, error) {
	// 先查找缓存
	if val, ok := optionCache.Load(id); ok {
		p = val.(*Option)
		return p, nil
	}
	// 缓存找不到从数据库中找
	o := orm.NewOrm()
	p.OptionId = id
	if err := o.Read(p); err != nil {
		return p, err
	}
	return p, nil
}

// FindByName 根据名称查找配置项
func (p *Option) FindByName(name string) (*Option, error) {
	// 先查找缓存
	if val, ok := optionCache.Load(name); ok {
		p = val.(*Option)
		return p, nil
	}
	// 缓存找不到从数据库中找
	o := orm.NewOrm()
	if err := o.QueryTable(p).Filter("option_name", name).One(p); err != nil {
		return p, err
	}
	return p, nil
}

// GetOptionValue 根据配置项名称获取配置项的值
func GetOptionValue(name, def string) string {
	if option, err := NewOption().FindByName(name); err == nil {
		return option.OptionValue
	}
	return def
}

// InsertOrUpdate 插入或更新配置项
func (p *Option) InsertOrUpdate() (err error) {
	defer func() {
		initOptionCache()
	}()
	o := orm.NewOrm()
	if p.OptionId > 0 || o.QueryTable(p.TableNameWithPrefix()).Filter("option_name", p.OptionName).Exist() {
		// 配置项存在则更新
		_, err = o.Update(p)
	} else {
		// 配置项不存在则插入
		_, err = o.Insert(p)
	}
	return err
}

// InsertMulti 批量插入配置项
func (p *Option) InsertMulti(option ...Option) error {
	o := orm.NewOrm()
	_, err := o.InsertMulti(len(option), option)
	initOptionCache()
	return err
}

// All 查找所有配置项
func (p *Option) All() ([]*Option, error) {
	o := orm.NewOrm()
	var options []*Option
	_, err := o.QueryTable(p.TableNameWithPrefix()).All(&options)
	if err != nil {
		return options, err
	}
	return options, nil
}

// ForbiddenReferer 禁止的Referer
func (m *Option) ForbiddenReferer() []string {
	return strings.Split(GetOptionValue("FORBIDDEN_REFERER", ""), "\n")
}
