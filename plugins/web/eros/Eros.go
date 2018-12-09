package eros

import (
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/utils"
	"github.com/cihub/seelog"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ecdiy/goserver/plugins"
)

func   Eros(ele *utils.Element) {
	au := &AppUpdater{
		Prefix:  ele.MustAttr("Prefix"),
		AppsDir: ele.MustAttr("AppDir")}
	au.initApp()
	Web :=plugins.GetGin(ele)
	//"/app/check"
	Web.GET(ele.Attr("CheckUrl", "/app/check"), au.Check)
	Web.GET(ele.Attr("ReloadUrl", "/app/reload"), au.Reload)

}

type AppUpdater struct {
	AppsDir, Prefix, LastVersion string
}

///app/check?jsVersion=72f32bf37c3377e205d1e6c06b47645e&appName=molove&android=0.0.1&isDiff=1
func (au *AppUpdater) Check(c *gin.Context) {
	wb := utils.NewParam(c)
	jsVersion := wb.String("jsVersion")
	android := wb.String("android")
	if android != "" {
		o := make(map[string]interface{})
		if au.LastVersion == jsVersion {
			o["resCode"] = 4000
			o["msg"] = "当前版本已是最新，不需要更新"
		} else {
			o["resCode"] = 0
			o["msg"] = "需要更新"
			o["Data"] = map[string]interface{}{
				"diff": false, "jsVersion": au.LastVersion,
				"path": au.Prefix + au.LastVersion + ".zip"}
		}
		c.JSON(200, o)
		return
	}
}

func (au *AppUpdater) Reload(c *gin.Context) {
	au.initApp()
	c.String(200, au.LastVersion)
}

func (au *AppUpdater) initApp() {
	au.LastVersion = ""
	file := au.AppsDir + "/bundle.config"
	conf, e := ioutil.ReadFile(file)
	if e == nil {
		m := make(map[string]interface{})
		e := json.Unmarshal(conf, &m)
		if e == nil {
			ver := fmt.Sprint(m["jsVersion"])
			au.LastVersion = ver
			seelog.Info("App Version:", ver)
			return
		} else {
			seelog.Error("", e.Error())
			return
		}
	} else {
		seelog.Error("EROS文件错误：", file, e)
		return
	}
}
