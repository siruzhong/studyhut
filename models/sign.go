package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"studyhut/constant"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	if _, err := os.Stat(constant.SignCacheDir); err != nil {
		err = os.MkdirAll(constant.SignCacheDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

// Sign 签到
type Sign struct {
	Id        int       // 主键id
	Uid       int       `orm:"index"` // 签到的用户id
	Day       int       `orm:"index"` // 签到日期，如20200101
	Reward    int       // 奖励的阅读秒数
	CreatedAt time.Time // 创建时间
}

// TableUnique 多字段唯一键
func (m *Sign) TableUnique() [][]string {
	return [][]string{
		[]string{"uid", "day"},
	}
}

// NewSign 新建签到实例
func NewSign() *Sign {
	return &Sign{}
}

// Rule 签到奖励规则
type Rule struct {
	BasicReward         int
	ContinuousReward    int
	MaxContinuousReward int
}

// 签到奖励规则全局对象
var _rule = &Rule{}

// GetSignRule 获取签到奖励规则
func (m *Sign) GetSignRule() (r *Rule) {
	return _rule
}

// UpdateSignRule 更新签到奖励规则
func (m *Sign) UpdateSignRule() {
	ops := []string{"SIGN_BASIC_REWARD", "SIGN_CONTINUOUS_REWARD", "SIGN_CONTINUOUS_MAX_REWARD"}
	for _, op := range ops {
		num, _ := strconv.Atoi(GetOptionValue(op, ""))
		switch op {
		case "SIGN_BASIC_REWARD":
			_rule.BasicReward = num
		case "SIGN_CONTINUOUS_REWARD":
			_rule.ContinuousReward = num
		case "SIGN_CONTINUOUS_MAX_REWARD":
			_rule.MaxContinuousReward = num
		}
	}
}

// IsSignToday 今天是否已签到
func (m *Sign) IsSignToday(uid int) bool {
	s := &Sign{}
	orm.NewOrm().QueryTable(m).Filter("uid", uid).Filter("day", time.Now().Format(constant.SignDayLayout)).One(s, "id")
	return s.Id > 0
}

// LatestSignTime 上次签到时间
func (m *Sign) LatestSignTime(uid int) (date int) {
	s := &Sign{}
	err := orm.NewOrm().QueryTable(m).Filter("uid", uid).OrderBy("-id").One(s)
	if err != nil {
		return
	}
	return int(s.CreatedAt.Unix())
}

// IsContinuousSign 是否未断签
func (m *Sign) IsContinuousSign(uid int) bool {
	s := &Sign{}
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour).Format(constant.SignDayLayout)
	orm.NewOrm().QueryTable(m).Filter("uid", uid).Filter("day", yesterday).One(s)
	return s.Id > 0
}

// Sign 执行签到（使用事务）
func (m *Sign) Sign(uid int) (reward int, err error) {
	s := &Sign{}
	o := orm.NewOrm()
	now := time.Now()
	day, _ := strconv.Atoi(now.Format(constant.SignDayLayout))
	// 1. 检测用户有没有签到
	o.QueryTable(s).Filter("uid", uid).Filter("day", day).One(s)
	if m.IsSignToday(uid) {
		err = errors.New(constant.MessageSigned)
		return
	}
	isContinuousSign := m.IsContinuousSign(uid) // 昨天有没有断签
	// 2. 查询用户签到了多少天
	user := NewMember()
	cols := []string{"member_id", "total_sign", "total_continuous_sign", "history_total_continuous_sign"}
	o.QueryTable(user).Filter("member_id", uid).One(user, cols...)
	if user.MemberId < 0 {
		err = errors.New(constant.MessageNotExistUser)
		return
	}
	// 3. 查询奖励规则
	rule := s.GetSignRule()
	// 4. 更新用户签到记录、签到天数和连续签到天数
	o.Begin()
	defer func() {
		if err != nil {
			beego.Error(err)
			err = errors.New(constant.MessageSignInnerErr)
			o.Rollback()
		} else {
			o.Commit()
		}
	}()
	user.TotalSign += 1
	s.Day = day
	s.Uid = uid
	s.CreatedAt = now
	//  奖励计算
	if isContinuousSign { // 连续签到
		user.TotalContinuousSign += 1
		extra := user.TotalContinuousSign * rule.ContinuousReward
		if extra >= rule.MaxContinuousReward {
			extra = rule.MaxContinuousReward
		}
		s.Reward = rule.BasicReward + extra
	} else { // 未连续签到
		user.TotalContinuousSign = 1
		s.Reward = rule.BasicReward + rule.ContinuousReward
	}
	if user.TotalContinuousSign > user.HistoryTotalContinuousSign {
		user.HistoryTotalContinuousSign = user.TotalContinuousSign
	}
	if _, err = o.Insert(s); err != nil {
		return
	}
	_, err = o.QueryTable(user).Filter("member_id", user.MemberId).Update(orm.Params{
		"total_sign":                    user.TotalSign,
		"total_continuous_sign":         user.TotalContinuousSign,
		"history_total_continuous_sign": user.HistoryTotalContinuousSign,
	})

	rt := NewReadingTime()
	o.QueryTable(rt).Filter("uid", uid).Filter("day", now.Format(constant.SignDayLayout)).One(rt)
	if rt.Id > 0 {
		rt.Duration += s.Reward
		_, err = o.Update(rt)
	} else {
		rt.Day, _ = strconv.Atoi(now.Format(constant.SignDayLayout))
		rt.Uid = uid
		rt.Duration = s.Reward
		_, err = o.Insert(rt)
	}
	reward = s.Reward
	return
}

// Sorted 签到排行
func (m *Sign) Sorted(limit int, orderField string, withCache ...bool) (members []Member) {
	var b []byte
	cache := false
	if len(withCache) > 0 {
		cache = withCache[0]
	}
	file := fmt.Sprintf(constant.SignCacheFmt, orderField, limit)
	if cache {
		if info, err := os.Stat(file); err == nil && time.Now().Sub(info.ModTime()).Seconds() <= cacheTime {
			// 文件存在，且在缓存时间内
			if b, err = ioutil.ReadFile(file); err == nil {
				json.Unmarshal(b, &members)
				if len(members) > 0 {
					return
				}
			}
		}
	}
	member := NewMember()
	o := orm.NewOrm()
	fields := []string{"member_id", "account", "nickname", "total_continuous_sign", "total_sign", "total_reading_time", "history_total_continuous_sign"}
	o.QueryTable(member).Filter("no_rank", 0).OrderBy("-"+orderField).Limit(limit).All(&members, fields...)
	if cache && len(members) > 0 {
		b, _ = json.Marshal(members)
		ioutil.WriteFile(file, b, os.ModePerm)
	}
	return
}

// LatestOne 上次签到
func (*Sign) LatestOne(uid int) (s Sign) {
	orm.NewOrm().QueryTable(&s).Filter("uid", uid).OrderBy("-id").One(&s)
	return
}

// SortedByPeriod 签到时间段排行
func (m *Sign) SortedByPeriod(limit int, prd period, withCache ...bool) (members []Member) {
	var b []byte
	cache := false
	if len(withCache) > 0 {
		cache = withCache[0]
	}
	file := fmt.Sprintf(constant.SignCacheFmt, "month-"+prd, limit)
	if cache {
		if info, err := os.Stat(file); err == nil && time.Now().Sub(info.ModTime()).Seconds() <= cacheTime {
			// 文件存在，且在缓存时间内
			if b, err = ioutil.ReadFile(file); err == nil {
				json.Unmarshal(b, &members)
				if len(members) > 0 {
					return
				}
			}
		}
	}
	sqlSort := "SELECT t.uid member_id,count(t.id) total_sign,m.account,m.avatar,m.nickname FROM `sign` t left JOIN members m on t.uid=m.member_id WHERE t.day>=? and t.day<=? GROUP BY t.uid ORDER BY total_sign desc limit ?"
	start, end := getTimeRange(time.Now(), prd)
	orm.NewOrm().Raw(sqlSort, start, end, limit).QueryRows(&members)
	if cache && len(members) > 0 {
		b, _ = json.Marshal(members)
		ioutil.WriteFile(file, b, os.ModePerm)
	}
	return
}
