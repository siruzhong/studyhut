package utils

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
)

var (
	AudioExt sync.Map // 音频扩展名
	VideoExt sync.Map // 视频拓展名
)

// init 初始化支持的音视频格式
func init() {
	// 音频格式
	for _, ext := range []string{".flac", ".wma", ".weba", ".aac", ".oga", ".ogg", ".mp3", ".webm", ".mid", ".wav", ".opus", ".m4a", ".amr", ".aiff", ".au"} {
		AudioExt.Store(ext, true)
	}
	// 视频格式
	for _, ext := range []string{".ogm", ".wmv", ".asx", ".mpg", ".webm", ".mp4", ".ogv", ".mpeg", ".mov", ".m4v", ".avi"} {
		VideoExt.Store(ext, true)
	}
}

// FileExists 判断指定的文件或目录是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer dst.Close()
	return io.Copy(dst, src)
}

func FormatBytes(size int64) string {
	units := []string{" B", " KB", " MB", " GB", " TB"}

	s := float64(size)

	i := 0

	for ; s >= 1024 && i < 4; i++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f%s", s, units[i])
}

func Round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

// ParseTitleFromMdHtml 从md的html文件中提取文章标题（从h1-h6）
func ParseTitleFromMdHtml(html string) (title string) {
	hTag := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(html)); err == nil {
		for _, tag := range hTag {
			if title = strings.TrimSpace(doc.Find(tag).First().Text()); title != "" {
				return title
			}
		}
	} else {
		beego.Error(err.Error())
	}
	return "空标题文档"
}

// ExecuteViewPathTemplate 执行指定的模板并返回执行结果.
//@param			tplName				模板文件路径
func ExecuteViewPathTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	viewPath := beego.BConfig.WebConfig.ViewsPath
	if err := beego.ExecuteViewPathTemplate(&buf, tplName, viewPath, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetUploadFileExt 获取允许的商城文件的类型
func GetUploadFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")

	temp := strings.Split(ext, "|")

	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

// IsAllowUploadFileExt 判断是否是允许商城的文件类型
func IsAllowUploadFileExt(ext string, typ ...string) bool {
	if len(typ) > 0 {
		t := strings.ToLower(strings.TrimSpace(typ[0]))
		if t == "audio" {
			_, ok := AudioExt.Load(ext)
			return ok
		} else if t == "video" {
			_, ok := VideoExt.Load(ext)
			return ok
		}
	}

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()

	for _, item := range exts {
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}
