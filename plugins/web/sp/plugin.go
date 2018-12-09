package sp

import (
	"github.com/ecdiy/goserver/plugins/web"
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins/sp"
)

func init() {

	//sp := &sp.WebSp{Gpa: plugins.GetGpa(ele) }
	//sp.Init()
	//doSubElement(ele, sp)
	//put(ele, sp)
	web.RegisterWebPlugin("SpHandle", HandleBase)
	web.RegisterWebPlugin("SpHandleCaptcha", HandleCaptcha)
}

type Handle struct {
	rule                  bool
	 spSuffix, RuleSp string
	ws                    *sp.WebSp
}

func (sh *Handle) Handle(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("sp un catch error;",   err)
		}
	}()
	spName := ctx.Param("sp") + sh.spSuffix
	wb := utils.NewParam(ctx)
	if sh.rule {
		ruleSp := sh.ws.GetRunSp(sh.RuleSp)
		if ruleSp == nil {
			seelog.Error("Rule Sp not exist.", sh.RuleSp)
			ctx.AbortWithStatus(403)
			return
		}
		params, rCode := sh.ws.GetParams(wb, ruleSp)
		if rCode != 200 {
			ctx.AbortWithStatus(403)
			return
		}
		r, _ := ruleSp.GetInt64(sh.ws.Gpa.Conn, params...)
		if r == 0 {
			ctx.AbortWithStatus(403)
			return
		}
	}

	code := sh.ws.SpExec(spName, wb)
	if code == 200 {
		ctx.JSON(200, wb.Out)
	} else {
		ctx.AbortWithStatus(code)
	}
}
