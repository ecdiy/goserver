package webs

import (
	"github.com/gin-gonic/gin"
	"utils/xml"
	"strings"
	"github.com/cihub/seelog"
	"strconv"
)

func (ws *SpWeb) Template(ele *xml.Element, data map[string]interface{}) {
	extends := make(map[string][]string)
	ns := ele.AllNodes()
	for _, n := range ns {
		extends[n.Name()] = strings.Split(n.Value, ",")
	}
	loginUrl := ele.MustAttr("LoginUrl")
	SpSuffix := ele.MustAttr("SpSuffix")
	ws.Engine.HTMLRender = WebTplRenderCreate(ele.MustAttr("TemplatesDir"), ele.MustAttr("LayoutDir"), extends, strings.Split(ele.MustAttr("Pages"), ",")...)
	rendFun := func(ctx *gin.Context) {
		if strings.Index(ctx.Request.URL.Path, ".") > 0 {
			seelog.Warn("404:", ctx.Request.URL.Path)
			ctx.AbortWithStatus(404)
			return
		}
		tplName := ws.GetTplName(ctx)
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
		if code == 200 || code == 404 {
			if code == 404 {
				seelog.Warn("Not Find SpName:", spName)
			}
			defer func() {
				if err := recover(); err != nil {
					seelog.Error("template error;template=", tplName+"-"+wb.Ua, "\nData=", wb.Out, "\n\n", err)
				}
			}()
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
	ws.Engine.NoRoute(rendFun)
}

func (ws *SpWeb) GetTplName(ctx *gin.Context) string {
	//if len(p.TplName) > 1 {
	//	return p.TplName
	//}
	url := ctx.Request.URL.Path
	if len(url) == 1 {
		url = "/index"
	}
	return url[1:]
}
