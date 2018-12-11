package QrCode

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/plugins/web"
	"github.com/cihub/seelog"
	"strconv"
)

func init() {
	web.RegisterWebPlugin("QrCode", func(ele *utils.Element) func(c *gin.Context) {
		w := ele.Attr("Width", "100")
		wInt, _ := strconv.Atoi(w)
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Cache-Control", "no-Cache, no-store, must-revalidate")
			c.Header("Pragma", "no-Cache")
			c.Header("Expires", "0")
			c.Header("Content-Type", "image/png")
			url := c.Query("url")
			qr, err := New(url, Medium)
			if err == nil {
				qr.Write(wInt, c.Writer)
			} else {
				seelog.Error("生成二维码出错")
			}
		}
	})
}
