package plugins

import (
	"github.com/ecdiy/goserver/utils"
	"os"
	"strings"
)

type BaseFun func(param *utils.Param, ps ... interface{}) interface{}

var Data = make(map[string]interface{}) //xml 对象保存

var pluginsMap = make(map[string]func(xml *utils.Element) interface{})

var ElementMap = make(map[string]*utils.Element)

var InitAfterFun []func() //xml 分析完后的回调函数

func RegisterPlugin(pluginName string, plugin func(ele *utils.Element) interface{}) {
	_, ext := pluginsMap[pluginName]
	if ext {
		panic("插件已存在：" + pluginName)
	} else {
		pluginsMap[pluginName] = plugin
	}
}

func GetRefByName(  name string) interface{} {

	dv, dvb := Data[name]
	if !dvb {
		panic("不存在:" + name)
	}
	return dv
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

func PutFunRun(fun func()) {
	if len(InitAfterFun) > 0 {
		go fun()
	} else {
		InitAfterFun = append(InitAfterFun, fun)
	}
}

func put(ele *utils.Element, v interface{}) {
	id, idb := ele.AttrValue("Id")
	if !idb {
		id = ele.Name()
	}
	_, de := Data[id]
	if de {
		panic("Id重复" + id)
	} else {
		Data[id] = v
	}
}

func getFile(file string) string {
	dir := os.Args[1]
	fg := []string{"/", "\\"}
	for _, flg := range fg {
		lst := strings.LastIndex(dir, flg)
		if lst > 0 {
			return dir[0:lst+1] + file
		}
	}
	return file
}
