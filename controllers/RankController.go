package controllers

import (
	"studyhut/constant"
	"studyhut/models"
)

// RankController 榜单控制器
type RankController struct {
	BaseController
}

// Index 榜单页面
func (this *RankController) Index() {
	limit, _ := this.GetInt("limit", 50)
	if limit > 200 {
		limit = 200
	}
	tab := this.GetString("tab", "all")
	switch tab {
	case "reading":
		rt := models.NewReadingTime()
		this.Data["SeoTitle"] = "阅读时长榜"
		this.Data["TodayReading"] = rt.Sort(constant.PeriodDay, limit, true)
		this.Data["WeekReading"] = rt.Sort(constant.PeriodWeek, limit, true)
		this.Data["MonthReading"] = rt.Sort(constant.PeriodMonth, limit, true)
		this.Data["LastWeekReading"] = rt.Sort(constant.PeriodLastWeek, limit, true)
		this.Data["LastMonthReading"] = rt.Sort(constant.PeriodLastMoth, limit, true)
		this.Data["AllReading"] = rt.Sort(constant.PeriodAll, limit, true)
	case "sign":
		this.Data["SeoTitle"] = "用户签到榜"
		sign := models.NewSign()
		this.Data["ContinuousSignUsers"] = sign.Sorted(limit, "total_continuous_sign", true)
		this.Data["TotalSignUsers"] = sign.Sorted(limit, "total_sign", true)
		this.Data["HistoryContinuousSignUsers"] = sign.Sorted(limit, "history_total_continuous_sign", true)
		this.Data["ThisMonthSign"] = sign.SortedByPeriod(limit, constant.PeriodMonth, true)
		this.Data["LastMonthSign"] = sign.SortedByPeriod(limit, constant.PeriodLastMoth, true)
	case "popular":
		this.Data["SeoTitle"] = "文档人气榜"
		bookCounter := models.NewBookCounter()
		//this.Data["Today"] = bookCounter.PageViewSort(constant.PeriodDay, limit, true)
		this.Data["Week"] = bookCounter.PageViewSort(constant.PeriodWeek, limit, true)
		this.Data["Month"] = bookCounter.PageViewSort(constant.PeriodMonth, limit, true)
		//this.Data["LastWeek"] = bookCounter.PageViewSort(constant.PeriodLastWeek, limit, true)
		//this.Data["LastMonth"] = bookCounter.PageViewSort(constant.PeriodLastMoth, limit, true)
		this.Data["All"] = bookCounter.PageViewSort(constant.PeriodAll, limit, true)
	case "star":
		this.Data["SeoTitle"] = "热门收藏榜"
		bookCounter := models.NewBookCounter()
		//this.Data["Today"] = bookCounter.StarSort(constant.PeriodDay, limit, true)
		this.Data["Week"] = bookCounter.StarSort(constant.PeriodWeek, limit, true)
		this.Data["Month"] = bookCounter.StarSort(constant.PeriodMonth, limit, true)
		//this.Data["LastWeek"] = bookCounter.StarSort(constant.PeriodLastWeek, limit, true)
		//this.Data["LastMonth"] = bookCounter.StarSort(constant.PeriodLastMoth, limit, true)
		this.Data["All"] = bookCounter.StarSort(constant.PeriodAll, limit, true)
	default:
		tab = "all"
		this.Data["SeoTitle"] = "总榜"
		limit = 10
		sign := models.NewSign()
		book := models.NewBook()
		this.Data["ContinuousSignUsers"] = sign.Sorted(limit, "total_continuous_sign", true)
		this.Data["ThisMonthSign"] = sign.SortedByPeriod(limit, constant.PeriodMonth, true)
		this.Data["TotalReadingUsers"] = sign.Sorted(limit, "total_reading_time", true)
		this.Data["StarBooks"] = book.Sorted(limit, "star")
		this.Data["VcntBooks"] = book.Sorted(limit, "vcnt")
		this.Data["CommentBooks"] = book.Sorted(limit, "cnt_comment")
	}
	this.Data["Tab"] = tab
	this.Data["IsRank"] = true
	this.TplName = "rank/index.html"
}
