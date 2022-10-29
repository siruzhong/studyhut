package store

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var ModelStoreCos = NewCos()

type Cos struct {
	Client *cos.Client
}

// NewCos 创建Cos连接客户端(官方文档:https://cloud.tencent.com/document/product/436/31215)
func NewCos() *Cos {
	bucketURL, _ := url.Parse(beego.AppConfig.String("cos::Domain"))
	serviceURL, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", beego.AppConfig.String("cos::Region")))
	baseUrl := &cos.BaseURL{BucketURL: bucketURL, ServiceURL: serviceURL}
	// 永久密钥
	client := &Cos{
		Client: cos.NewClient(baseUrl, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  beego.AppConfig.String("cos::SecretId"),
				SecretKey: beego.AppConfig.String("cos::SecretKey"),
			},
		}),
	}
	return client
}

// GetBucket 获取bucket
func (c *Cos) GetBucket() cos.Bucket {
	s, _, err := NewCos().Client.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}
	return s.Buckets[0]
}

// IsObjectExist 判断文件对象是否存在
func (c *Cos) IsObjectExist(object string) (err error) {
	ok, err := NewCos().Client.Object.IsExist(context.Background(), object)
	if err == nil && ok {
		fmt.Printf("object exists\n")
	} else if err != nil {
		fmt.Printf("head object failed: %v\n", err)
	} else {
		fmt.Printf("object does not exist\n")
	}
	return
}

// MoveToCos 文件移动到Cos进行存储
// @param	local		本地文件
// @param	save		存储到Cos的文件
// @param	IsDel		文件上传后，是否删除本地文件
// @param	IsGzip		是否做gzip压缩，做gzip压缩的话，需要修改oss中对象的响应头，设置gzip响应
func (c *Cos) MoveToCos(local, save string, IsDel bool, IsGzip ...bool) error {
	isGzip := false
	// 如果是开启了gzip，则需要设置文件对象的响应头
	if len(IsGzip) > 0 && IsGzip[0] == true {
		isGzip = true
	}
	// 在移动文件到Cos之前，先压缩文件
	if isGzip {
		if bs, err := ioutil.ReadFile(local); err != nil {
			beego.Error(err.Error())
			isGzip = false // 设置为false
		} else {
			var by bytes.Buffer
			w := gzip.NewWriter(&by)
			defer w.Close()
			w.Write(bs)
			w.Flush()
			ioutil.WriteFile(local, by.Bytes(), 0777)
		}
	}
	_, _, err := NewCos().Client.Object.Upload(context.Background(), save, local, nil)
	if err != nil {
		beego.Error("文件移动到Cos失败：", err.Error())
		return err
	}
	// 删除上传后的本地文件
	if err == nil && IsDel {
		err = os.Remove(local)
	}
	return err
}

// DelFromCos 从CoS中删除文件
// @param	object		文件对象
func (c *Cos) DelFromCos(object ...string) (err error) {
	var objects []string
	objects = append(objects, object...)
	obs := []cos.Object{}
	for _, v := range objects {
		obs = append(obs, cos.Object{Key: v})
	}
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
	}
	_, _, err = NewCos().Client.Object.DeleteMulti(context.Background(), opt)
	return err
}

// DelCosFolder 根据cos文件夹
func (o *Cos) DelCosFolder(folder string) (err error) {
	_, err = NewCos().Client.Object.Delete(context.Background(), folder)
	return
}

//// HandleContent 处理html中的Cos数据：如果是用于预览的内容，则把img等的链接的相对路径转成绝对路径，否则反之
//// @param	htmlstr		html字符串
//// @param	forPreview	是否是供浏览的页面需求
//// @return	str			处理后返回的字符串
//func (c *Cos) HandleContent(htmlStr string, forPreview bool) (str string) {
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
//	if err != nil {
//		beego.Error(err.Error())
//		return htmlStr
//	}
//	doc.Find("img").Each(func(i int, s *goquery.Selection) {
//		if src, exist := s.Attr("src"); exist {
//			//预览
//			if forPreview {
//				// 存在http开头的图片链接，则更新为绝对链接
//				if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
//					s.SetAttr("src", beego.AppConfig.String("cos::Domain")+"/"+strings.TrimLeft(src, "./"))
//				}
//			} else {
//				s.SetAttr("src", strings.TrimPrefix(src, beego.AppConfig.String("cos::Domain")))
//			}
//		}
//	})
//	str, _ = doc.Find("body").Html()
//	return
//}
//
//// DelByHtmlPics 从HTML中提取图片文件，并删除
//func (c *Cos) DelByHtmlPics(htmlStr string) {
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
//	if err != nil {
//		beego.Error(err.Error())
//		return
//	}
//	doc.Find("img").Each(func(i int, s *goquery.Selection) {
//		// For each item found, get the band and title
//		if src, exist := s.Attr("src"); exist {
//			//不存在http开头的图片链接，则更新为绝对链接
//			if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
//				c.DelFromCos(strings.TrimLeft(src, "./")) //删除
//			} else if strings.HasPrefix(src, beego.AppConfig.String("cos::Domain")) {
//				c.DelFromCos(strings.TrimPrefix(src, beego.AppConfig.String("cos::Domain"))) //删除
//			}
//		}
//	})
//	return
//}
