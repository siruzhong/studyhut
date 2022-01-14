package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Local 本地存储
type Local struct{}

var ModelStoreLocal = new(Local)

// IsObjectExist 判断文件对象是否存在(err为nil表示文件存在，否则表示文件不存在，并告知错误信息)
func (this *Local) IsObjectExist(object string) (err error) {
	_, err = os.Stat(object)
	return
}

// MoveToStore 文件存储
// @param	tmpFile	临时文件
// @param	save	存储文件，不建议与临时文件相同
func (this *Local) MoveToStore(tmpFile, save string) (err error) {
	save = strings.TrimLeft(save, "/")
	//"./a.png"与"a.png"是相同路径
	if strings.HasPrefix(tmpFile, "./") || strings.HasPrefix(save, "./") {
		tmpFile = strings.TrimPrefix(tmpFile, "./")
		save = strings.TrimPrefix(save, "./")
	}
	if strings.ToLower(tmpFile) != strings.ToLower(save) { // 不是相同文件路径
		os.MkdirAll(filepath.Dir(save), os.ModePerm)
		// 不使用rename，因为在docker中会挂载外部卷，导致错误
		// 见https://gocn.vip/article/178
		if b, err := ioutil.ReadFile(tmpFile); err == nil {
			ioutil.WriteFile(save, b, os.ModePerm)
		}
		os.Remove(tmpFile)
	}
	return
}

// DelFiles 本地删除文件
func (this *Local) DelFiles(object ...string) error {
	for _, file := range object {
		os.Remove(strings.TrimLeft(file, "/"))
	}
	return nil
}

// DelFolder 本地删除文件夹
func (this *Local) DelFolder(folder string) (err error) {
	return os.RemoveAll(folder)
}
