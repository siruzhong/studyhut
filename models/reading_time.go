package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"studyhut/constant"

	"github.com/astaxie/beego/orm"
)

// ReadingTime 阅读时长
type ReadingTime struct {
	Id       int
	Uid      int
	Day      int // 日期，如 20191212
	Duration int // 每天的阅读时长
}

// sum 总阅读时长
type sum struct {
	SumVal int
}

type ReadingSortedUser struct {
	Uid              int    `json:"uid"`
	Account          string `json:"account"`
	Nickname         string `json:"nickname"`
	Avatar           string `json:"avatar"`
	SumTime          int    `json:"sum_time"`
	TotalReadingTime int    `json:"total_reading_time"`
}

const (
	readingTimeCacheDir = "cache/rank/reading-time"
	readingTimeCacheFmt = "cache/rank/reading-time/%v-%v.json"
)

func init() {
	if _, err := os.Stat(readingTimeCacheDir); err != nil {
		err = os.MkdirAll(readingTimeCacheDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

// NewReadingTime 创建阅读时长对象
func NewReadingTime() *ReadingTime {
	return &ReadingTime{}
}

func (*ReadingTime) TableUnique() [][]string {
	return [][]string{[]string{"uid", "day"}}
}

// GetReadingTime 获取阅读时长
func (r *ReadingTime) GetReadingTime(uid int, prd period) int {
	sum := &sum{}
	o := orm.NewOrm()
	sqlSum := "select sum(duration) sum_val from reading_time where uid = ? and day>=? and day<=? limit 1"
	now := time.Now()
	if prd == constant.PeriodAll {
		m := NewMember()
		o.QueryTable(m).Filter("member_id", uid).One(m, "total_reading_time")
		return m.TotalReadingTime
	}
	start, end := getTimeRange(now, prd)
	o.Raw(sqlSum, uid, start, end).QueryRow(sum)
	return sum.SumVal
}

func (r *ReadingTime) Sort(prd period, limit int, withCache ...bool) (users []ReadingSortedUser) {
	var b []byte
	cache := false
	if len(withCache) > 0 {
		cache = withCache[0]
	}
	file := fmt.Sprintf(readingTimeCacheFmt, prd, limit)
	if cache {
		if info, err := os.Stat(file); err == nil && time.Now().Sub(info.ModTime()).Seconds() <= cacheTime {
			// 文件存在，且在缓存时间内
			if b, err = ioutil.ReadFile(file); err == nil {
				json.Unmarshal(b, &users)
				if len(users) > 0 {
					return
				}
			}
		}
	}

	sqlSort := "SELECT t.uid,sum(t.duration) sum_time,m.account,m.avatar,m.nickname FROM `reading_time` t left JOIN members m on t.uid=m.member_id WHERE m.no_rank=0 and t.day>=? and t.day<=? GROUP BY t.uid ORDER BY sum_time desc limit ?"
	start, end := getTimeRange(time.Now(), prd)
	orm.NewOrm().Raw(sqlSort, start, end, limit).QueryRows(&users)

	if cache && len(users) > 0 {
		b, _ = json.Marshal(users)
		ioutil.WriteFile(file, b, os.ModePerm)
	}
	return
}
