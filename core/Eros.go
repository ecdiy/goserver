package core

import (
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/utils"
	"github.com/cihub/seelog"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/ecdiy/goserver/webs"
)

var (
	AppsDir    string
	Prefix     string
	AppVersion = make(map[string]string)
	Apps       = []string{"mo"}
)

func (app *Module) Eros(ele *utils.Element) {
	AppsDir = utils.EnvParam("Apps")
	Prefix = utils.EnvParam("Prefix")
	seelog.Info("Apps", Apps)

	Web := getGin(ele)
	Web.GET("/app/check", appCheck)
	Web.GET("/app/reload", appReload)
	initApp()
}

func initApp() string {
	res := ""
	for _, app := range Apps {
		initAppByName(app)
	}
	info := "\n\t\tAppsDir=" + AppsDir + ";Prefix=" + Prefix +
		";AppVersion=" + fmt.Sprint(AppVersion) + ";\nReloadRes=" + res
	seelog.Info(info)
	return info
}

func initAppByName(app string) string {
	file := AppsDir + "/" + app + "/bundle.config"
	conf, e := ioutil.ReadFile(file)
	if e == nil {
		m := make(map[string]interface{})
		e := json.Unmarshal(conf, &m)
		if e == nil {
			ver := fmt.Sprint(m["jsVersion"])

			zipF := AppsDir + "/" + ver + ".zip"
			_, e := os.Stat(zipF)
			if os.IsNotExist(e) {
				seelog.Info("!!not exist.", zipF)
				bf, _ := ioutil.ReadFile(AppsDir + "/" + app + "/bundle.zip")
				ioutil.WriteFile(zipF, bf, 0644)
			} else {
				seelog.Info("exist zip file.", zipF)
			}

			AppVersion[app] = ver
			return file + ";\npath=" + Prefix + "/apps/" + ver + ".zip;\n"
		} else {
			seelog.Error("", e.Error())
			return ""
		}
	} else {
		return ",not exist=" + fmt.Sprint(e)
	}
}

func appReload(c *gin.Context) {
	c.String(200, initApp())
}

///app/check?jsVersion=72f32bf37c3377e205d1e6c06b47645e&appName=molove&android=0.0.1&isDiff=1
func appCheck(c *gin.Context) {
	wb := webs.NewParam(c)
	jsVersion := wb.String("jsVersion")
	appName := wb.String("appName")
	android := wb.String("android")
	if android != "" {
		mv, mvb := AppVersion[appName]
		if !mvb || len(AppVersion) == 0 {
			seelog.Info("没有缓存,", appName)
			initAppByName(appName)
			mv, mvb = AppVersion[appName]
		}
		o := make(map[string]interface{})
		if mvb && mv == jsVersion {
			o["resCode"] = 4000
			o["msg"] = "当前版本已是最新，不需要更新"
		} else {
			seelog.Info("update...,feq Ver=", jsVersion)
			o["resCode"] = 0
			o["msg"] = "需要更新"
			o["data"] = map[string]interface{}{
				"diff":      false,
				"jsVersion": mv,
				"path":      Prefix + "/apps/" + mv + ".zip"}
		}
		c.JSON(200, o)
		return
	}
	//TODO...
	seelog.Info("appName=", appName)
}
