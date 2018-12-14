package resize

import (
	"github.com/ecdiy/goserver/plugins/web"
	"github.com/gin-gonic/gin"
	"strings"
	"strconv"
	"io/ioutil"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins/web/image"
)

func init() {
	web.RegisterWebPlugin("Resize", func(ele *utils.Element) func(c *gin.Context) {
		ir := &ImageResize{dir: ele.MustAttr("Dir"),
			size: strings.Split(ele.MustAttr("Size"), ",")}
		return ir.DoResize
	})
}

type ImageResize struct {
	dir  string
	size []string
}

func (ir *ImageResize) DoResize(context *gin.Context) {
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
							image.ImgResize(source, target, tw)
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
