package sp

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"github.com/ecdiy/goserver/plugins/sp"
)

func HandleCaptcha(ele *utils.Element) func(c *gin.Context) {
	ws := &sp.WebSp{}
	spName := ele.MustAttr("Sp")
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("sp un catch error;", err)
			}
		}()
		wb := utils.NewParam(ctx)
		if !captcha.VerifyString(wb.String("CaptchaId"), wb.String("CaptchaVal")) {
			wb.ST(utils.StErrorCaptcha, captcha.New())
			ctx.JSON(200, wb.Out)
			return
		}
		code := ws.SpExec(spName, wb)
		if code == 200 {
			ctx.JSON(200, wb.Out)
		} else {
			seelog.Error("数据存储过程错误:"+spName, ";", code)
			ctx.AbortWithStatus(code)
		}
	}
}
