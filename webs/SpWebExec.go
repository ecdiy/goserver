package webs

import (
	"utils/gpa"
	"github.com/cihub/seelog"
	"utils"
)

/**
URL 映射到存储过程调用，返回json数据格式
SpAjax(true,"/sp/","Sp",authFun)
 */

//func WebSp(Gpa *gpa.Gpa, eng *gin.Engine, auth func(c *gin.Context) (bool, int64), adminAuthFun func(c *gin.Context) (bool, int64), reload func()) {
//	if !gin.IsDebugging() {
//		if adminAuthFun != nil {
//			spInitCache(Gpa, adminAuthFun, "Admin")
//		}
//		if auth != nil {
//			spInitCache(Gpa, auth, "Ajax")
//			spInitCache(Gpa, auth, "Page")
//		}
//	}
//	if auth != nil {
//		eng.POST("/sp/:sp", func(c *gin.Context) {
//			defer func() {
//				if err := recover(); err != nil {
//					seelog.Error("sp un catch error;", err)
//				}
//			}()
//			sp(Gpa, c, "Ajax", auth)
//		})
//	}
//	if adminAuthFun != nil {
//		eng.POST("/spa/:sp", func(c *gin.Context) {
//			defer func() {
//				if err := recover(); err != nil {
//					seelog.Error("spa un catch error;", err)
//				}
//			}()
//			sp(Gpa, c, "Admin", adminAuthFun)
//		})
//	}
//	eng.GET("/spReload", func(c *gin.Context) {
//		defer func() {
//			if err := recover(); err != nil {
//				seelog.Error("spReload un catch error;", err)
//			}
//		}()
//		seelog.Info("Reload Sp Cache & template ...")
//		if reload != nil {
//			reload()
//		}
//		spCache = make(map[string]*Sp)
//		for fix, fun := range spReloadFun {
//			spInitCache(Gpa, fun, fix)
//		}
//		utils.OK.OutJSON(c, nil)
//	})
//}


func SpExec(spName string, g *gpa.Gpa, param *Param, auth func(c *Param) *UserBase) int {
	defer func() {
		if err := recover(); err != nil {
			delete(spCache, spName)
		}
	}()
	var sp *Sp
	var ext bool
	if utils.EnvIsDev {
		sp, ext = NewSpByName(g, spName, auth)
	} else {
		sp, ext = spCache[spName]
		if !ext {
			sp, ext = NewSpByName(g, spName, auth)
			if ext {
				spCache[spName] = sp
			}
		}
	}
	if !ext {
		return 404
	}
	params, code := sp.GetParams(param)
	if code == 200 {
		e := sp.Run(param.Out, g.Conn, params...)
		if e != nil {
			seelog.Error("exec SP失败:", sp.Name)
			delete(spCache, sp.Name)
			return 500
		}
		return 200
	} else {
		return code
	}
}
