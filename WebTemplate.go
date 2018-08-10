package utils

import (
	"github.com/cihub/seelog"
	"github.com/gin-contrib/multitemplate"
	"path/filepath"
	"strings"
)

// WebTplRenderCreate("ui/views/wk-site/", "/**/*", "/*")
func WebTplRenderCreate(templatesDir string, pages ...string) multitemplate.Renderer {
	ext := ".html"
	r := multitemplate.NewRenderer()
	webMs, _ := filepath.Glob(templatesDir + "layout/web/*" + ext)
	h5Ms, _ := filepath.Glob(templatesDir + "layout/h5/*" + ext)
	seelog.Info("web modules:", webMs)
	for _, p := range pages {
		pages, err := filepath.Glob(templatesDir + p + ext)
		if err != nil {
			panic(err.Error())
		}
		for _, page := range pages {
			j := strings.LastIndex(page, "-")
			if j < 0 || j < len(templatesDir) {
				if strings.Index(page, "layout") < 0 {
					seelog.Warn("命名规则不对:-ua[].html")
				}
				continue
			}
			layout := templatesDir + "layout/" + page[j+1:]
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
			seelog.Info(n, ms)
			r.AddFromFiles(n, ms...)
		}
	}
	return r
}
