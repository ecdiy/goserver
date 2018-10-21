package webs

import (
	"github.com/gin-gonic/gin"
	"goserver/utils"
	"strings"
	"github.com/cihub/seelog"
	"strconv"
)

func (ws *SpWeb) Template(ele *utils.Element, data map[string]interface{}) {
	extends := make(map[string][]string)
	ns := ele.Node("extends")
	if ns != nil {
		for _, n := range ns.AllNodes() {
			extends[n.MustAttr("TplName")] = strings.Split(n.Value, ",")
		}
	}
	loginUrl := ele.MustAttr("LoginUrl")
	SpSuffix := ele.MustAttr("SpSuffix")
	ws.Engine.HTMLRender = WebTplRenderCreate(ele.MustAttr("TemplatesDir"), ele.MustAttr("LayoutDir"), extends, strings.Split(ele.MustAttr("Pages"), ",")...)
	mp := ele.Node("map")
	if mp != nil {
		its := mp.AllNodes()
		for _, it := range its {
			ws.Engine.GET(it.MustAttr("Url"), ws.getTemplateRender(data, SpSuffix, loginUrl, it.MustAttr("TplName")))
		}
	}
	ws.Engine.NoRoute(ws.getTemplateRender(data, SpSuffix, loginUrl, ""))
}

func (ws *SpWeb) getTemplateRender(data map[string]interface{}, SpSuffix, loginUrl, TplName string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		defer func() {
			//if err := recover(); err != nil {
			//	seelog.Error("模板处理失败", err)
			//}
		}()
		if strings.Index(ctx.Request.URL.Path, ".") > 0 {
			seelog.Warn("404:", ctx.Request.URL.Path)
			ctx.AbortWithStatus(404)
			return
		}
		tplName := TplName
		if len(TplName) == 0 {
			tplName = ws.GetTplName(ctx)
		}
		ns := strings.Split(tplName, "/")
		spName := ""
		for _, n := range ns {
			if len(n) > 1 {
				spName += strings.ToUpper(n[0:1]) + n[1:]
			}
		}
		spName += SpSuffix
		wb := NewParam(ctx)
		code := ws.SpExec(spName, wb)

		wb.Out["gin"] = wb

		if code == 200 || code == 404 {
			ctx.HTML(200, tplName+"-"+wb.Ua, wb.Out)
		} else {
			if code == 401 {
				ctx.Redirect(302, loginUrl)
			} else {
				seelog.Error("code=", code, ",spName=", spName, ",tplName=", tplName)
				ctx.HTML(200, strconv.Itoa(code)+"-"+wb.Ua, wb)
			}
		}
	}
}

func (ws *SpWeb) GetTplName(ctx *gin.Context) string {
	url := ctx.Request.URL.Path
	if len(url) == 1 {
		url = "/index"
	}
	return url[1:]
}
