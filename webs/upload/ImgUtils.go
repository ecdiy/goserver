package upload

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/cihub/seelog"
	"github.com/hunterhug/go_image/graphics"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

func FmtImgDir(dir, uri string) (string, string) {
	ln := len(uri)
	if ln > 6 {
		x := uri[0:3] + "/" + uri[3:6]
		os.MkdirAll(dir+x, 0644)
		urx := x + "/" + uri[6:]
		path := dir + urx
		return path, urx
	}
	if ln > 3 {
		x := uri[0:3]
		os.MkdirAll(dir+x, 0644)
		urx := x + "/" + uri[3:]
		path := dir + urx
		return path, urx
	}
	return dir + uri, uri
}

func Md5Byte(bs []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(bs)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func Md5File(fileName string) (string, error) {
	bs, e := ioutil.ReadFile(fileName)
	if e == nil {
		return Md5Byte(bs), nil
	} else {
		return "", e
	}
}

func loadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

func ImgResize(sourceFile, outFile string, newWidth int) (int, int) {
	src, err := loadImage(sourceFile)
	if err != nil {
		seelog.Error("图片缩放失败:", sourceFile, " ==> ", outFile, " ~~ ", err)
		return 0, 0
	}
	bound := src.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newWidth*dy/dx))
	err = graphics.Scale(dst, src)
	if err != nil {
		seelog.Error(err)
	}
	out, err := os.Create(outFile)
	if err != nil {
		seelog.Error(err)
	}
	defer out.Close()
	if strings.LastIndex(strings.ToLower(sourceFile), ".png") > 0 {
		png.Encode(out, dst)
	} else {
		jpeg.Encode(out, dst, nil)
	}
	return dx, dy
}
