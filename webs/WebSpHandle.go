package webs

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"strings"
	"github.com/ecdiy/goserver/utils"
)

func (ws *WebSp) HandleCaptcha(ele *utils.Element, data map[string]interface{}) {
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

func (ws *WebSp) Handle(ele *utils.Element, data map[string]interface{}) {
	sh := &SpHandle{url: ele.MustAttr("Url"), ws: ws, spSuffix: ele.MustAttr("SpSuffix")}

	if !gin.IsDebugging() {
		list, err := ws.Gpa.ListArrayString(SqlSpAll)
		if err != nil {
			seelog.Error("查询所有的存储过程出错：", err)
		} else {
			if list != nil {
				sps := ""
				for _, val := range list {
					if strings.LastIndex(val[0], sh.spSuffix) == len(val[0])-len(sh.spSuffix) {
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

	sh.RuleSp, sh.rule = ele.AttrValue("RuleSp")

	ws.Engine.POST(sh.url, sh.Handle)
	reloadUrl, rExt := ele.AttrValue("ReloadUrl")
	if rExt {
		ws.Engine.GET(reloadUrl, func(i *gin.Context) {
			ws.SpCache = make(map[string]*Sp)
			i.String(200, "clear cache ok.")
		})
	}
}

type SpHandle struct {
	rule                  bool
	url, spSuffix, RuleSp string
	ws                    *WebSp
}

func (sh *SpHandle) Handle(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("sp un catch error;", sh.url, ";", err)
		}
	}()
	spName := ctx.Param("sp") + sh.spSuffix
	wb := NewParam(ctx)
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
		r, _ := ruleSp.GetInt64(sh.ws.Gpa.Conn, params)
		if r == 0 {
			ctx.AbortWithStatus(403)
			return
		}
	}

	code := sh.ws.SpExec(spName, wb)
	if code == 200 {
		ctx.JSON(200, wb.Out)
	} else {
		seelog.Error("数据存储过程错误:"+spName, ";", code)
		ctx.AbortWithStatus(code)
	}
}
