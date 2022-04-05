package models

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TruthHun/gotil/sitemap"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"studyhut/constant"
)

// 时间段
type period string

var cacheTime = beego.AppConfig.DefaultFloat("CacheTime", 60) // 1 分钟

var (
	AllowRegister = true
	AllowVisitor  = true
)

func Init() {
	initAPI()
	initAdsCache()
	initOptionCache()
	NewSign().UpdateSignRule()          // 更新签到规则的全局变量
	NewReadRecord().UpdateReadingRule() // 更新阅读计时规则的全局变量
	go func() {
		for {
			AllowRegister = GetOptionValue("ENABLED_REGISTER", "true") == "true"
			AllowVisitor = GetOptionValue("ENABLE_ANONYMOUS", "true") == "true"
			time.Sleep(time.Second * 30)
		}
	}()
}

// SetIncreAndDecre 设置增减
//@param            table           需要处理的数据表
//@param            field           字段
//@param            condition       条件
//@param            incre           是否是增长值，true则增加，false则减少
//@param            step            增或减的步长
func SetIncreAndDecre(table string, field string, condition string, incre bool, step ...int) (err error) {
	mark := "-"
	if incre {
		mark = "+"
	}
	s := 1
	if len(step) > 0 {
		s = step[0]
	}
	sql := fmt.Sprintf("update %v set %v=%v%v%v where %v", table, field, field, mark, s, condition)
	_, err = orm.NewOrm().Raw(sql).Exec()
	return
}

type SitemapDocs struct {
	DocumentId   int
	DocumentName string
	Identify     string
	BookId       int
}

// SitemapData 站点地图数据
func SitemapData(page, listRows int) (totalRows int64, sitemaps []SitemapDocs) {
	//获取公开的书籍
	var (
		books   []Book
		docs    []Document
		maps    = make(map[int]string)
		booksId []interface{}
	)

	o := orm.NewOrm()
	o.QueryTable("books").Filter("privately_owned", 0).Limit(100000).All(&books, "book_id", "identify")
	if len(books) > 0 {
		for _, book := range books {
			booksId = append(booksId, book.BookId)
			maps[book.BookId] = book.Identify
		}
		q := o.QueryTable("documents").Filter("BookId__in", booksId...)
		totalRows, _ = q.Count()
		q.Limit(listRows).Offset((page-1)*listRows).All(&docs, "document_id", "document_name", "book_id")
		if len(docs) > 0 {
			for _, doc := range docs {
				sd := SitemapDocs{
					DocumentId:   doc.DocumentId,
					DocumentName: doc.DocumentName,
					BookId:       doc.BookId,
				}
				if v, ok := maps[doc.BookId]; ok {
					sd.Identify = v
				}
				sitemaps = append(sitemaps, sd)
			}
		}
	}
	return
}

// SitemapUpdate 更新站点地图
func SitemapUpdate(domain string) {
	var (
		files   []string
		bookIds []interface{}
		bookMap = make(map[int]string)
		Sitemap = sitemap.NewSitemap("1.0", "utf-8")
		o       = orm.NewOrm()
		si      []sitemap.SitemapIndex
	)
	domain = strings.TrimSuffix(domain, "/")
	os.Mkdir("sitemap", os.ModePerm)
	//查询公开的书籍
	qsBooks := o.QueryTable("books").Filter("privately_owned", 0)
	limit := 10000
	for i := 0; i < 10; i++ {
		var books []Book
		qsBooks.Limit(limit).Offset(i*limit).All(&books, "book_id", "identify", "release_time", "book_name")
		if len(books) > 0 {
			file := "sitemap/books-" + strconv.Itoa(i) + ".xml"
			files = append(files, file)
			var su []sitemap.SitemapUrl
			for _, book := range books {
				su = append(su, sitemap.SitemapUrl{
					Loc:        domain + beego.URLFor("DocumentController.Index", ":key", book.Identify),
					Lastmod:    book.ReleaseTime.Format("2006-01-02 15:04:05"),
					ChangeFreq: sitemap.WEEKLY,
					Priority:   0.9,
				})
				bookIds = append(bookIds, book.BookId)
				bookMap[book.BookId] = book.Identify
			}
			Sitemap.CreateSitemapContent(su, file)
		} else {
			i = 10
		}
	}
	qsDocs := o.QueryTable("documents").Filter("book_id__in", bookIds...)
	for i := 0; i < 100; i++ {
		var docs []Document
		qsDocs.Limit(limit).Offset(i*limit).All(&docs, "modify_time", "book_id", "document_name", "document_id", "identify")
		if len(docs) > 0 {
			file := "sitemap/docs-" + strconv.Itoa(i) + ".xml"
			files = append(files, file)
			var su []sitemap.SitemapUrl
			for _, doc := range docs {
				bookIdentify := ""
				if idtf, ok := bookMap[doc.BookId]; ok {
					bookIdentify = idtf
				}
				su = append(su, sitemap.SitemapUrl{
					Loc:        domain + beego.URLFor("DocumentController.Read", ":key", bookIdentify, ":id", url.QueryEscape(doc.Identify)),
					Lastmod:    doc.ModifyTime.Format("2006-01-02 15:04:05"),
					ChangeFreq: sitemap.WEEKLY,
					Priority:   0.9,
				})
			}
			Sitemap.CreateSitemapContent(su, file)
		} else {
			i = 100
		}
	}
	if len(files) > 0 {
		for _, f := range files {
			si = append(si, sitemap.SitemapIndex{
				Loc:     domain + "/" + f,
				Lastmod: time.Now().Format("2006-01-02 15:04:05"),
			})
		}

	}
	Sitemap.CreateSitemapIndex(si, "sitemap.xml")
}

type Count struct {
	Cnt        int
	CategoryId int
}

// CountCategory 统计书籍分类
func CountCategory() {
	var count []Count

	o := orm.NewOrm()
	sql := "select count(bc.id) cnt, bc.category_id from book_category bc left join books b on b.book_id=bc.book_id where b.privately_owned=0 and bc.category_id>0  group by bc.category_id"
	o.Raw(sql).QueryRows(&count)
	if len(count) == 0 {
		return
	}

	var cates []Category
	tableCate := "category"
	o.QueryTable(tableCate).All(&cates, "id", "pid", "cnt")
	if len(cates) == 0 {
		return
	}

	var err error

	o.Begin()
	defer func() {
		if err != nil {
			beego.Error(err)
			o.Rollback()
		} else {
			o.Commit()
		}
	}()

	cateChild := make(map[int]int)
	if _, err = o.QueryTable(tableCate).Filter("id__gt", 0).Update(orm.Params{"cnt": 0}); err != nil {
		return
	}

	for _, item := range count {
		cateChild[item.CategoryId] = item.Cnt
		_, err = o.QueryTable(tableCate).Filter("id", item.CategoryId).Update(orm.Params{"cnt": item.Cnt})
		if err != nil {
			return
		}
	}
}

// getTimeRange 获取时间范围
func getTimeRange(t time.Time, prd period) (start, end string) {
	switch prd {
	case constant.PeriodWeek:
		start, end = getWeek(t)
	case constant.PeriodLastWeek:
		start, end = getWeek(t.AddDate(0, 0, -7))
	case constant.PeriodMonth:
		start, end = getMonth(t)
	case constant.PeriodLastMoth:
		start, end = getMonth(t.AddDate(0, -1, 0))
	case constant.PeriodAll:
		start = "20060102"
		end = "20401231"
	case constant.PeriodDay:
		start = t.Format(constant.DateFormat)
		end = start
	case constant.PeriodYear:
		start, end = getYear(t.AddDate(-1, 0, 0))
	default:
		start = t.Format(constant.DateFormat)
		end = start
	}
	return
}

// getWeek 获取周时间
func getWeek(t time.Time) (start, end string) {
	if t.Weekday() == 0 {
		start = t.Add(-7 * 24 * time.Hour).Format(constant.DateFormat)
		end = t.Format(constant.DateFormat)
	} else {
		s := t.Add(-time.Duration(t.Weekday()-1) * 24 * time.Hour)
		start = s.Format(constant.DateFormat)
		end = s.Add(6 * 24 * time.Hour).Format(constant.DateFormat)
	}
	return
}

// getYear 获取年时间
func getYear(t time.Time) (start, end string) {
	month := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	start = month.Format(constant.DateFormat)
	end = month.AddDate(0, 12, 0).Add(-24 * time.Hour).Format(constant.DateFormat)
	return
}

// getMonth 获取月时间
func getMonth(t time.Time) (start, end string) {
	month := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	start = month.Format(constant.DateFormat)
	end = month.AddDate(0, 1, 0).Add(-24 * time.Hour).Format(constant.DateFormat)
	return
}
