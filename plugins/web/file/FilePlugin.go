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
	path := p.Dir + wb.String("path")
	if action == "mkdir" {
		os.MkdirAll(path, 0644)
	}
	if action == "delFile" {
		os.Remove(path)
		return
	}
	if action == "md2Html" {
		bs, e := ioutil.ReadFile(path)
		if e == nil {
			c.JSON(200, Md2Html(bs))
		} else {
			c.JSON(200, "")
		}
		return
	}

	if action == "read" {
		bs, e := ioutil.ReadFile(path)
		if e == nil {
			c.Data(200, wb.String("contentType"), bs)
		} else {
			c.String(200, "")
		}
		return
	}

	if action == "save" {
		ix := strings.LastIndex(path, "/")
		if ix < 0 {
			ix = strings.LastIndex(path, "\\")
		}
		if ix > 0 {
			dir := path[0:ix]
			os.MkdirAll(dir, 0644)
		}
		out, err := os.Create(path)
		if err == nil {
			out.Write([]byte(wb.String("body")))
		}
		out.Close()
		return
	}
}
