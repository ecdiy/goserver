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

func SpAjax(uri string, g *gpa.Gpa, eng *gin.Engine, spPrefix string, auth func(c *gin.Context) (bool, int64)) {
	eng.GET(uri+"Reload", func(c *gin.Context) {
		seelog.Info("Reload Sp Cache")
		spCache = make(map[string]*Sp)
		spInitCache(g, auth)
		utils.OK.OutJSON(c, nil)
	})
	if !gin.IsDebugging() {
		spInitCache(g, auth)
	}
	eng.POST(uri+"/*sp", func(c *gin.Context) {
		spName := spPrefix + c.Param("sp")
		var sp *Sp
		var ext bool
		if gin.IsDebugging() {
			sp, ext = NewSpByName(g, spName, auth)
		} else {
			sp, ext = spCache[spName]
		}
		if !ext {
			seelog.Error("sp not exist.", sp)
			c.AbortWithStatus(404)
			return
		}
		wb := WebBaseNew(c)
		data, code := SpExec(g, sp, wb)
		if code == 200 {
			c.JSON(200, data)
		} else {
			seelog.Error("数据存储过程错误:"+spName, "\n\t", code)
			c.AbortWithStatus(code)
		}
	})
}

func SpExec(g *gpa.Gpa, sp *Sp, ctx *WebBase) (map[string]interface{}, int) {
	defer func() {
		if err := recover(); err != nil {
			delete(spCache, sp.Name)
		}
	}()
	params, code := sp.GetParams(ctx)
	if code == 200 {
		m, e := sp.Run(g, params...)
		if e != nil {
			seelog.Error("exec SP失败:", sp.Name)
			delete(spCache, sp.Name)
			return nil, 500
		}
		return m, 200
	} else {
		return nil, code
	}
}
