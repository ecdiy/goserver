package web

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/plugins"
	"strings"
)

var webPluginsMap = make(map[string]func(xml *utils.Element) func(c *gin.Context))
var pluginsMap = make(map[string]func(xml *utils.Element) func(c *utils.Param))

func init() {
	plugins.RegisterPlugin("Web", func(ele *utils.Element) interface{} {
		web := gin.New()
		ns := ele.AllNodes()
		for _, n := range ns {
			if n.Name() == "Static" {
				web.Static(n.MustAttr("Url"), n.MustAttr("Path"))
			} else {
				mth := strings.ToUpper(n.Attr("Method", "POST"))
				url, ub := n.AttrValue("Url")
				if !ub {
					panic("Web Attr设置错误:" + n.Name() + ";缺少Url")
				}
				pp, ppFd := pluginsMap[n.Name()]
				if ppFd {
					imp := pp(n)
					f := func(c *gin.Context) {
						wb := utils.NewParam(c)
						imp(wb)
						c.JSON(200, wb.Out)
					}
					if mth == "GET" {
						web.GET(url, f)
					} else {
						web.POST(url, f)
					}
					continue
				}
				p, pFd := webPluginsMap[n.Name()]
				if pFd {
					if mth == "GET" {
						web.GET(n.MustAttr("Url"), p(n))
					} else {
						web.POST(n.MustAttr("Url"), p(n))
					}
				} else {

					panic("没有实现的Web方法:" + n.Name())
				}
			}
		}
		plugins.PutFunRun(func() {
			port, _ := ele.AttrValue("Port")
			web.Run(port)
		})
		return web
	})
}

func RegisterWebPlugin(pluginName string, plugin func(ele *utils.Element) func(c *gin.Context)) {
	webPluginsMap[pluginName] = plugin
}
func RegisterWeb(pluginName string, plugin func(ele *utils.Element) func(c *utils.Param)) {
	pluginsMap[pluginName] = plugin
}