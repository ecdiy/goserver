package main

import (
	"os"
	"github.com/cihub/seelog"
	"github.com/gpmgo/gopm/modules/log"
	"reflect"
	"utils/xml"
)

var app = reflect.ValueOf(new(Module))
var data = make(map[string]interface{}) //xml 对象保存
var initAfterFun []func()               //xml 分析完后的回调函数

func main() {
	defer func() {
		seelog.Flush()
		if err := recover(); err != nil {
			seelog.Error("sp un catch error;", err)
			panic(err)
		}
	}()
	if len(os.Args) < 2 {
		seelog.Error("没有指定配置文件")
		return
	}
	xmlFile := "./conf/" + os.Args[1] + ".xml"
	dom, err := xml.LoadByFile(xmlFile)
	if err == nil {
		InvokeByXml(dom)
		if len(initAfterFun) > 0 {
			for _, fun := range initAfterFun {
				fun()
			}
		}
	} else {
		seelog.Error("读取配置文件出错:", xmlFile, ",", err)
		return
	}
}

func putFunRun(fun func()) {
	if len(initAfterFun) > 0 {
		go fun()
	} else {
		initAfterFun = append(initAfterFun, fun)
	}
}

func put(ele *xml.Element, v interface{}) {
	id, idb := ele.AttrValue("Id")
	if !idb {
		id = ele.Name()
	}
	_, de := data[id]
	if de {
		panic("Id重复" + id)
	} else {
		data[id] = v
		seelog.Info("~~~map object:", id)
	}
}

func InvokeByXml(ecXml *xml.Element) {
	ns := ""
	defer func() {
		seelog.Info("analysis element:", ns, "\n\tdata:", data)
	}()
	allNode := ecXml.AllNodes()
	log.Debug("配置节点数:", len(allNode))
	for _, n := range allNode {
		ns += n.Name() + ";"
		//if IsReload {
		//	canReload, canReloadBool := n.AttrValue("canReload")
		//	if !canReloadBool || canReload != "1" {
		//		continue
		//	}
		//	log.Info("~~~reload~~~", n.Name())
		//}
		inputs := make([]reflect.Value, 1)
		inputs[0] = reflect.ValueOf(n)
		m := app.MethodByName(n.Name())
		if m.IsValid() {
			m.Call(inputs)
		} else {
			panic("没有实现的方法:" + n.Name())
		}
	}
}
