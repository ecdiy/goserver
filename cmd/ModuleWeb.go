package main

import (
	"goserver/webs"
	"github.com/gin-gonic/gin"
	"goserver/utils"
	"goserver/webs/upload"
)

func (app *Module) Web(ele *utils.Element) {
	web := gin.New()
	ns := ele.AllNodes()
	for _, n := range ns {
		if n.Name() == "Static" {
			web.Static(n.MustAttr("Url"), n.MustAttr("Path"))
		}
	}
	put(ele, web)
	putFunRun(func() {
		port, _ := ele.AttrValue("Port")
		web.Run(port)
	})
}

func (app *Module) Sp(ele *utils.Element) {
	sp := &webs.WebSp{Gpa: getGpa(ele), Engine: getGin(ele)}
	sp.Init()
	doSubElement(ele, sp)
	put(ele, sp)
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
func (app *Module) Upload(ele *utils.Element) {
	sp, spExt := ele.AttrValue("SpRef")
	if spExt {
		upload.Upload(getGin(ele), data[sp].(*webs.WebSp), getVerify(ele), ele)
	} else {
		upload.Upload(getGin(ele), nil, getVerify(ele), ele)
	}
}
