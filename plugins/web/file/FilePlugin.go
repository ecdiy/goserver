package file

import (
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/webs"
	"os"
	"io/ioutil"
	"strings"
	"github.com/ecdiy/goserver/plugins"
)

func init() {
	plugins.WebPlugins["File"] = func(xml *utils.Element) func(c *gin.Context) {
		p := &Plugin{Dir: xml.MustAttr("Dir")}
		return p.Impl
	}
}

type Plugin struct {
	Dir string
}

func (p *Plugin) Impl(c *gin.Context) {
	wb := webs.NewParam(c)
	action := wb.String("action")
	path := wb.String("path")
	if action == "mkdir" {
		os.MkdirAll(p.Dir+path, 0644)
	}
	if action == "md2Html" {
		bs, e := ioutil.ReadFile(p.Dir + path)
		if e == nil {
			c.JSON(200, Md2Html(bs))
		} else {
			c.JSON(200, "")
		}
	}

	if action == "save" {
		pt := p.Dir + path
		ix := strings.LastIndex(pt, "/")
		if ix < 0 {
			ix = strings.LastIndex(pt, "\\")
		}
		if ix > 0 {
			dir := pt[0:ix]
			os.MkdirAll(dir, 0644)
		}
		out, err := os.Create(pt)
		if err == nil {
			out.Write([]byte(wb.String("body")))
		}
		out.Close()
	}
}
