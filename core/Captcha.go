package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/dchest/captcha"
	"github.com/ecdiy/goserver/webs"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/plugins"
)

/**
认证码
 */
func (app *Module) WebCaptcha(ele *utils.Element) {
	plugins.GetGin(ele).GET(ele.MustAttr("Url"), func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-Cache, no-store, must-revalidate")
		c.Header("Pragma", "no-Cache")
		c.Header("Expires", "0")
		c.Header("Content-Type", "image/png")
		id := c.Query("t")
		captcha.WriteImage(c.Writer, id, captcha.StdWidth, captcha.StdHeight)
	})
}

func (app *Module) WebCaptchaNew(ele *utils.Element) {
	post(ele, func(param *webs.Param) {
		param.OK(captcha.New())
	})
}
