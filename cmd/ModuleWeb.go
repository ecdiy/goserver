package main

import (
	"utils/webs"
	"github.com/gin-gonic/gin"
	"utils/xml"
)

func (app *Module) Web(ele *xml.Element) {
	web := gin.New()
	put(ele, web)
	putFunRun(func() {
		port, _ := ele.AttrValue("Port")
		web.Run(port)
	})
}

func (app *Module) Sp(ele *xml.Element) {
	sp := &webs.SpWeb{Gpa: getGpa(ele), Engine: getGin(ele)}
	sp.Init()
	doSubElement(ele, sp)
	put(ele, sp)
}
