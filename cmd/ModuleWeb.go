package main

import (
	"utils/webs"
	"github.com/gin-gonic/gin"
	"utils/xml"
	"reflect"
	"github.com/cihub/seelog"
)

func (app *Module) Web(ele *xml.Element) {
	web := gin.New()
	put(ele, web)
	putFunRun(func() {
		port, _ := ele.AttrValue("Port")
		web.Run(port)
	})
}

func (app *Module) WebSp(ele *xml.Element) {
	sp := &webs.SpWeb{Gpa: getGpa(ele), SpSuffix: ele.MustAttr("SpSuffix"), SpParamDoMap: map[string]webs.ParamValFunc{}}
	ns := ele.AllNodes()
	if len(ns) > 0 {
		spv := reflect.ValueOf(sp)
		for _, e := range ns {
			inputs := make([]reflect.Value, 2)
			inputs[0] = reflect.ValueOf(e)
			inputs[1] = reflect.ValueOf(data)
			m := spv.MethodByName(e.Name())
			seelog.Warn("sp register :" + e.Name())
			m.Call(inputs)
		}
	}
	sp.Handle(getGin(ele), ele.MustAttr("Url"), ele.MustAttr("ReloadUrl"))
	put(ele, sp)
}
