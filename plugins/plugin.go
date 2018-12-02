package plugins

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/gpa"
)

var Data = make(map[string]interface{}) //xml 对象保存

var Plugins = make(map[string]func(xml *utils.Element))

var WebPlugins = make(map[string]func(xml *utils.Element) func(c *gin.Context))

var ElementMap = make(map[string]*utils.Element)

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