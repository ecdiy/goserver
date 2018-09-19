package webs

import (
	"github.com/gin-gonic/gin"
	"utils/gpa"
	"github.com/cihub/seelog"
	"utils"
)

/**
URL 映射到存储过程调用，返回json数据格式
SpAjax(true,"/sp/","Sp",authFun)
 */

func WebSp(g *gpa.Gpa, eng *gin.Engine, auth func(c *gin.Context) (bool, int64), adminAuthFun func(c *gin.Context) (bool, int64), reload func()) {
	if !gin.IsDebugging() {
		if adminAuthFun != nil {
			spInitCache(g, adminAuthFun, "Admin")
		}
		if auth != nil {
			spInitCache(g, auth, "Ajax")
			spInitCache(g, auth, "Page")
		}
	}
	if auth != nil {
		eng.POST("/sp/:sp", func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					seelog.Error("sp un catch error;", err)
				}
			}()
			sp(g, c, "Ajax", auth)
		})
	}
	if adminAuthFun != nil {
		eng.POST("/spa/:sp", func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					seelog.Error("spa un catch error;", err)
				}
			}()
			sp(g, c, "Admin", adminAuthFun)
		})
	}
	eng.GET("/spReload", func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Error("spReload un catch error;", err)
			}
		}()
		seelog.Info("Reload Sp Cache & template ...")
		if reload != nil {
			reload()
		}
		spCache = make(map[string]*Sp)
		for fix, fun := range spReloadFun {
			spInitCache(g, fun, fix)
		}
		utils.OK.OutJSON(c, nil)
	})
}

func sp(g *gpa.Gpa, c *gin.Context, spPrefix string, auth func(c *gin.Context) (bool, int64)) {
	spName := c.Param("sp") + spPrefix
	wb := WebBaseNew(c)
	code := SpExec(spName, g, wb, auth)
	if code == 200 {
		c.JSON(200, wb.Out)
	} else {
		seelog.Error("数据存储过程错误:"+spName, ";", code)
		c.AbortWithStatus(code)
	}
}

func SpExec(spName string, g *gpa.Gpa, ctx *WebBase, auth func(c *gin.Context) (bool, int64)) int {
	defer func() {
		if err := recover(); err != nil {
			delete(spCache, spName)
		}
	}()
	var sp *Sp
	var ext bool
	if gin.IsDebugging() {
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
	params, code := sp.GetParams(ctx)
	if code == 200 {
		e := sp.Run(ctx.Out, g.Conn, params...)
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
