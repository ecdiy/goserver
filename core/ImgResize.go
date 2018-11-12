package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/ecdiy/goserver/webs/upload"
	"strconv"
	"io/ioutil"
	"github.com/cihub/seelog"
)

type ImgResize struct {
	dir  string
	size []string
}

func (ir *ImgResize) DoResize(context *gin.Context) {
	url := context.Query("url")
	if len(url) > 1 {
		dp := strings.LastIndex(url, ".")
		pp := strings.LastIndex(url, "_")
		if dp > pp && pp > 0 {
			w := url[pp+1 : dp]
			for _, ss := range ir.size {
				if w == ss {
					target := ir.dir + url
					source := ir.dir + url[0:pp] + url[dp:]
					if strings.Index(w, "x") < 0 {
						tw, twe := strconv.Atoi(w)
						if twe == nil {
							upload.ImgResize(source, target, tw)
							bs, err := ioutil.ReadFile(target)
							if err == nil {
								context.Header("Content-Type", "image/"+url[dp+1:])
								context.Writer.Write(bs)
							} else {
								seelog.Error("target file error:", target)
							}
							return
						}
					}
				}
			}
			seelog.Error("生成的宽度不支持,", w)
		}
	}
}

func (app *Module) ImgResize(ele *utils.Element) {
	ir := &ImgResize{dir: ele.MustAttr("Dir"),
		size: strings.Split(ele.MustAttr("Size"), ",")}
	getGin(ele).GET(ele.MustAttr("Url"), ir.DoResize)
}
