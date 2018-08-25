package webs

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"utils"
	"utils/gpa"
)

type Sp struct {
	Sql, Name, SessionName string
	Params                 []*SpParam
	Result                 []*SpResult
	//Info                   map[string]interface{}
}

type SpResult struct {
	Name, Type string //[ [total,object],[list,list] ]
}
type SpParam struct {
	ParamName string
	ValFunc   ParamValFunc
}

type ParamValFunc func(ctx *gin.Context, p *SpParam) (interface{}, int)

func SpAjax(uri string, g *gpa.Gpa, eng *gin.Engine, spPrefix string, auth func(c *gin.Context) bool) {
	eng.GET(uri+"Reload", func(c *gin.Context) {
		seelog.Info("Reload Sp Cache")
		spCache = make(map[string]*Sp)
		spInitCache(g)
		utils.OK.OutJSON(c, nil)
	})
	if !gin.IsDebugging() {
		spInitCache(g)
	}
	if auth != nil {
		eng.POST(uri+"/*sp", func(c *gin.Context) {
			spName := spPrefix + c.Param("sp")
			var sp *Sp
			var ext bool
			if gin.IsDebugging() {
				sp, ext = NewSpByName(g, spName)
			} else {
				sp, ext = spCache[spName]
			}
			if !ext {
				seelog.Error("sp not exist.", sp)
				c.AbortWithStatus(404)
				return
			}
			data, code := Exec(g, sp, c)
			if code == 200 {
				c.JSON(200, data)
			} else {
				seelog.Error("数据存储过程错误:"+spName, "\n\t", code)
				c.AbortWithStatus(code)
			}
		})
	}
}
