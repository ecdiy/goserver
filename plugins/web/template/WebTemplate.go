package template

import (
	"github.com/cihub/seelog"
	"strings"
	"github.com/gin-gonic/gin"
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
	seelog.Info("web modules:", layoutDir, webMs, h5Ms)
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
