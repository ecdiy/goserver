package webs

import (
	"github.com/cihub/seelog"
	"strings"
	"github.com/gin-gonic/gin"
	"utils/gpa"
	"github.com/gin-contrib/multitemplate"
	"path/filepath"
	"io/ioutil"
	"regexp"
)


// WebTplRenderCreate("ui/views/wk-site/", "/**/*", "/*")
func WebTplRenderCreate(templatesDir, layoutDir string, extends map[string][]string, pages ...string) multitemplate.Renderer {

	regs := []*regexp.Regexp{regexp.MustCompile(">\\s*<"),
		regexp.MustCompile(";\\s*"), regexp.MustCompile(",\\s*"),
		regexp.MustCompile("{\\s*"), regexp.MustCompile("\\s+[}]\\s*"),
		regexp.MustCompile("[\r\n]+\\s*")}
	tos := []string{"><", ";", ",", "{", "}", "\n"}

	ext := ".html"
	r := multitemplate.NewRenderer()
	webMs, _ := filepath.Glob(layoutDir + "/web/*" + ext)
	h5Ms, _ := filepath.Glob(layoutDir + "/h5/*" + ext)
	seelog.Info("web modules:", webMs)
	for _, p := range pages {
		pp, err := filepath.Glob(templatesDir + p + ext)
		if err != nil {
			panic(err.Error())
		}
		for _, page := range pp {
			j := strings.LastIndex(page, "-")
			if j < 0 || j < len(templatesDir) {
				if strings.Index(page, "layout") < 0 {
					seelog.Warn("命名规则不对:-ua[].html,", page)
				}
				continue
			}
			layout := layoutDir + "/" + page[j+1:]
			n := page[0:len(page)-len(ext)][len(templatesDir):]
			d := strings.Index(n, ",")
			if d > 0 {
				n = n[0:d]
			}
			n = strings.Replace(n, "\\", "/", -1)
			ms := []string{layout, page}
			if strings.Index(n, "-web") > 0 {
				if len(webMs) > 0 {
					ms = append(ms, webMs...)
				}
			} else {
				if len(h5Ms) > 0 {
					ms = append(ms, h5Ms...)
				}
			}
			mds, mb := extends[n]
			if mb {
				ms = append(ms, mds...)
			}
			seelog.Info(n, ms)
			if gin.IsDebugging() {
				r.AddFromFilesFuncs(n, FunConstantMaps, ms...)
			} else {
				var html []string
				for _, mm := range ms {
					bb, _ := ioutil.ReadFile(mm)
					s := strings.TrimSpace(string(bb))
					for i, r := range regs {
						s = r.ReplaceAllString(s, tos[i])
					}
					html = append(html, s)
				}
				r.AddFromStringsFuncs(n, FunConstantMaps, html...)
			}
		}
	}
	return r
}

type WebTemplate struct {
	* Param
	TplName string //模版名称
}

func (p *WebTemplate) GetTplName() string {
	if len(p.TplName) > 1 {
		return p.TplName
	}
	url := p.Context.Request.URL.Path
	if len(url) == 1 {
		url = "/index"
	}
	return url[1:]
}

func WebTplWithSp(loginUrl string, tpl *WebTemplate, ctx *gin.Context, g *gpa.Gpa, auth func(c *Param)*UserBase) {
	if strings.Index(ctx.Request.URL.Path, ".") > 0 {
		seelog.Warn("404:", ctx.Request.URL.Path)
		ctx.AbortWithStatus(404)
		return
	}
	tplName := tpl.GetTplName()
	ns := strings.Split(tplName, "/")
	spName := ""
	for _, n := range ns {
		if len(n) > 1 {
			spName += strings.ToUpper(n[0:1]) + n[1:]
		}
	}
	spName += "Page"
	//code := SpExec(spName, g, tpl. Param, auth)
	//if code == 200 || code == 404 {
	//	if code == 404 {
	//		seelog.Warn("Not Find SpName:", spName)
	//	}
	//	defer func() {
	//		if err := recover(); err != nil {
	//			seelog.Error("template error;template=", tplName+"-"+tpl.Ua,
	//				"\nData=", tpl.Out, "\n\n", err)
	//		}
	//	}()
	//	ctx.HTML(200, tplName+"-"+tpl.Ua, tpl.Out)
	//} else {
	//	if code == 401 {
	//		ctx.Redirect(302, loginUrl)
	//	} else {
	//		seelog.Error("code=", code, ",spName=", spName, ",tplName=", tplName)
	//		ctx.HTML(200, strconv.Itoa(code)+"-"+tpl.Ua, tpl)
	//	}
	//}
}
