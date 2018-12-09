package plugins

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/gpa"
	"strings"
)

var Data = make(map[string]interface{}) //xml 对象保存

var Plugins = make(map[string]func(xml *utils.Element))

var WebPlugins = make(map[string]func(xml *utils.Element) func(c *gin.Context))

var ElementMap = make(map[string]*utils.Element)

var InitAfterFun []func() //xml 分析完后的回调函数

func GetGpa(ele *utils.Element) *gpa.Gpa {
	ref := ele.Attr("GpaRef", "Gpa")
	web := Data[ref].(*gpa.Gpa)
	return web
}

func GetGin(ele *utils.Element) *gin.Engine {
	ref := ele.Attr("WebRef", "Web")
	web := Data[ref].(*gin.Engine)
	return web
}

func Invoke(n *utils.Element) bool {
	w, we := WebPlugins[n.Name()]
	if we {
		mtd := strings.ToUpper(n.Attr("Method", "Get"))
		if strings.Index(mtd, "GET") >= 0 {
			GetGin(n).GET(n.MustAttr("Url"), w(n))
		} else {
			GetGin(n).POST(n.MustAttr("Url"), w(n))
		}
		return true
	}
	p, pFd := Plugins[n.Name()]
	if pFd {
		p(n)
		return true
	}
	return false
}
