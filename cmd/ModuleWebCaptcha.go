package main

import (
	"github.com/dchest/captcha"
	"utils/xml"
	"utils/webs"
	"github.com/gin-gonic/gin"
)

func (app *Module) WebCaptcha(ele *xml.Element) {
	getGin(ele).GET(ele.MustAttr("Url"), func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-Cache, no-store, must-revalidate")
		c.Header("Pragma", "no-Cache")
		c.Header("Expires", "0")
		c.Header("Content-Type", "image/png")
		id := c.Query("t")
		captcha.WriteImage(c.Writer, id, captcha.StdWidth, captcha.StdHeight)
	})
}

func (app *Module) WebCaptchaNew(ele *xml.Element) {
	post(ele, func(param *webs.Param) {
		param.OK(captcha.New())
	})
}
