package main

import (
	"utils/xml"
	"github.com/gin-gonic/gin"
	"utils/gpa"
)

func getGpa(ele *xml.Element) *gpa.Gpa {
	ref, _ := ele.AttrValue("GpaRef")
	web := data[ref].(*gpa.Gpa)
	return web
}
func getGin(ele *xml.Element) *gin.Engine {
	ref, _ := ele.AttrValue("WebRef")
	web := data[ref].(*gin.Engine)
	return web
}
