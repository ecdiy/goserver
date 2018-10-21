package webs

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"strings"
	"goserver/utils"
)

func (ws *SpWeb) HandleCaptcha(ele *utils.Element, data map[string]interface{}) {
	url := ele.MustAttr("Url")
	spName := ele.MustAttr("Sp")
	ws.Engine.POST(url, func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("sp un catch error;", url, ";", err)
			}
		}()
		wb := NewParam(ctx)
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
	})
}

func (ws *SpWeb) Handle(ele *utils.Element, data map[string]interface{}) {
	spSuffix := ele.MustAttr("SpSuffix")
	if !gin.IsDebugging() {
		list, err := ws.Gpa.ListArrayString(SqlSpAll)
		if err != nil {
			seelog.Error("查询所有的存储过程出错：", err)
		} else {
			if list != nil {
				sps := ""
				for _, val := range list {
					if strings.LastIndex(val[0], spSuffix) == len(val[0])-len(spSuffix) {
						sp, b := ws.NewSp(val)
						if b {
							ws.SpCache [sp.Name] = sp
							sps += sp.Sql + ","
						}
					}
				}
				seelog.Info("~~~~;\n\t", sps)
			}
		}
	}
	url := ele.MustAttr("Url")
	ws.Engine.POST(url, func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("sp un catch error;", url, ";", err)
			}
		}()
		spName := ctx.Param("sp") + spSuffix
		wb := NewParam(ctx)
		code := ws.SpExec(spName, wb)
		if code == 200 {
			ctx.JSON(200, wb.Out)
		} else {
			seelog.Error("数据存储过程错误:"+spName, ";", code)
			ctx.AbortWithStatus(code)
		}
	})
	reloadUrl, rExt := ele.AttrValue("ReloadUrl")
	if rExt {
		ws.Engine.GET(reloadUrl, func(i *gin.Context) {
			ws.SpCache = make(map[string]*Sp)
			i.String(200, "clear cache ok.")
		})
	}
}
