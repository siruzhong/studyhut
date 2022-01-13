package utils

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// ImageCopy 图片拷贝
func ImageCopy(src image.Image, x, y, w, h int) (image.Image, error) {
	var subImg image.Image
	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {
		return subImg, errors.New("图片解码失败")
	}
	return subImg, nil
}

// ImageCopyFromFile 从文件拷贝图片
func ImageCopyFromFile(p string, x, y, w, h int) (src image.Image, err error) {
	file, err := os.Open(p)
	if err != nil {
		return src, err
	}
	defer file.Close()
	src, _, err = image.Decode(file)
	return ImageCopy(src, x, y, w, h)
}

// ImageResize 调整图片大小
func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}

// ImageResizeSaveFile 调整图片大小并保存
func ImageResizeSaveFile(src image.Image, width, height int, p string) error {
	dst := resize.Resize(uint(width), uint(height), src, resize.Lanczos3)
	return SaveImage(p, dst)
}

// SaveImage 将图片保存到指定的路径
func SaveImage(p string, src image.Image) error {
	os.MkdirAll(filepath.Dir(p), 0666)
	f, err := os.OpenFile(p, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(p)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		err = jpeg.Encode(f, src, &jpeg.Options{Quality: 80})
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, src, &gif.Options{NumColors: 256})
	}
	return err
}
