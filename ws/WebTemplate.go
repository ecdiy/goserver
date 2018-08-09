package ws

import (
	"github.com/cihub/seelog"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

// WebCreateTemplateRender("ui/views/wk-site/", "/**/*", "/*")
func WebCreateTemplateRender(templatesDir string, pages ...string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	for _, p := range pages {
		ext := ".html"
		pages, err := filepath.Glob(templatesDir + p + ext)
		if err != nil {
			panic(err.Error())
		}
		for _, include := range pages {
			fb := filepath.Base(include)
			layout := templatesDir
			if strings.Index(fb, "-h5") > 0 {
				layout += "h5.htm"
			} else if strings.Index(fb, "-web") > 0 {
				layout += "web.htm"
			} else {
				seelog.Warn("no ua.", fb)
				continue
			}
			n := include[0 : len(include)-len(ext)][len(templatesDir):]
			n = strings.Replace(n, "\\", "/", -1)
			r.AddFromFiles(n, layout, include)
			seelog.Info(n, ":[", layout, ",", include, "]")
		}
	}
	return r
}

func WebHtml(url string, fun func(wdb *Web)) {
	WebGin.GET(url, func(c *gin.Context) {
		web := WebNew(c)
		fun(web)
		u := url
		if len(u) == 1 {
			u = "/index"
		}
		c.HTML(200, u[1:]+"-"+web.Ua, web)
	})
}
