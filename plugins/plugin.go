package plugins

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
)

type BaseFun func(param *utils.Param, ps ... interface{}) interface{}

var Data = make(map[string]interface{}) //xml 对象保存

var pluginsMap = make(map[string]func(xml *utils.Element) interface{})

var ElementMap = make(map[string]*utils.Element)

var InitAfterFun []func() //xml 分析完后的回调函数

func GetGin(ele *utils.Element) *gin.Engine {
	ref := ele.Attr("WebRef", "Web")
	web := Data[ref].(*gin.Engine)
	return web
}

func RegisterPlugin(pluginName string, plugin func(xml *utils.Element) interface{}) {
	_, ext := pluginsMap[pluginName]
	if ext {
		panic("插件已存在：" + pluginName)
	} else {
		pluginsMap[pluginName] = plugin
	}
}

func GetRef(ele *utils.Element, DefaultRef string) interface{} {
	bfId, vb := ele.AttrValue(DefaultRef + "Ref")
	if !vb {
		bfId = DefaultRef
	}
	dv, dvb := Data[bfId]
	if !dvb {
		panic("不存在:" + bfId)
	}
	return dv
}
