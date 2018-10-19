package main

import (
	"utils/xml"
	"github.com/gin-gonic/gin"
	"utils/gpa"
	"reflect"
	"github.com/cihub/seelog"
	"utils/webs"
)

func getGpa(ele *xml.Element) *gpa.Gpa {
	ref := ele.Attr("GpaRef", "Gpa")
	web := data[ref].(*gpa.Gpa)
	return web
}

func getGin(ele *xml.Element) *gin.Engine {
	ref := ele.Attr("WebRef", "Web")
	web := data[ref].(*gin.Engine)
	return web
}

func doSubElement(ele *xml.Element, obj interface{}) {
	ns := ele.AllNodes()
	if len(ns) > 0 {
		spv := reflect.ValueOf(obj)
		for _, e := range ns {
			inputs := make([]reflect.Value, 2)
			inputs[0] = reflect.ValueOf(e)
			inputs[1] = reflect.ValueOf(data)
			m := spv.MethodByName(e.Name())
			seelog.Warn("sp register :" + e.Name())
			m.Call(inputs)
		}
	}
}

func post(ele *xml.Element, fun func(param *webs.Param)) {
	getGin(ele).POST(ele.MustAttr("Url"), func(c *gin.Context) {
		wb := webs.NewParam(c)
		fun(wb)
		c.JSON(200, wb.Out)
	})
}
