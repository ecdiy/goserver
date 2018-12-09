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
				mth := strings.ToUpper(n.Attr("Method", "GET"))
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

func RegisterWebPlugin(pluginName string, plugin func(xml *utils.Element) func(c *gin.Context)) {
	webPluginsMap[pluginName] = plugin
}
func RegisterWeb(pluginName string, plugin func(xml *utils.Element) func(c *utils.Param)) {
	pluginsMap[pluginName] = plugin
}

//文件上传
/**
TmpDir,MainWidth,ImgWidth（多个用,分隔） 可选

	两种方式:
    1.上传完成后不调用存储过程，输出 文件名对应的参数Url
    <Upload WebRef="Web" TmpDir="./upload/temp/" DirUpload="./upload/" ImgWidth="800" MainWidth="800" UrlPrefix="/upload"/>

	2.上传完成后调用存储过程， 存储过程返回值决定输出内容
    <Upload SpRef="Sp" WebRef="Web" TmpDir="./upload/temp/" DirUpload="./upload/" ImgWidth="800" Sp="Upload" MainWidth="800" UrlPrefix="/upload"/>

 */

//func post(ele *utils.Element, fun func(param *utils.Param)) {
//	GetGin(ele).POST(ele.MustAttr("Url"), func(c *gin.Context) {
//		wb := utils.NewParam(c)
//		fun(wb)
//		c.JSON(200, wb.Out)
//	})
//}
