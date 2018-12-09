package captcha

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/plugins/web"
)

func init() { /**
认证码
 */
	web.RegisterWebPlugin("Captcha", func(xml *utils.Element) func(c *gin.Context) {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Cache-Control", "no-Cache, no-store, must-revalidate")
			c.Header("Pragma", "no-Cache")
			c.Header("Expires", "0")
			c.Header("Content-Type", "image/png")
			id := c.Query("t")
			captcha.WriteImage(c.Writer, id, captcha.StdWidth, captcha.StdHeight)
		}
	})
	web.RegisterWeb("CaptchaNew", func(xml *utils.Element) func(wb *utils.Param) {
		return func(wb *utils.Param) {
			wb.OK(captcha.New())
		}
	})
}
