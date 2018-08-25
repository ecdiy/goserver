package webs

import (
	"github.com/gin-gonic/gin"
	"utils/gpa"
)

/**
URL 映射到存储过程调用，返回json数据格式
SpAjax(true,"/sp/","Sp",authFun)
 */
func SpAjax(uri string, g *gpa.Gpa, eng *gin.Engine, spPrefix string, auth func(c *gin.Context) bool) {
	//eng.GET(uri+"Reload", func(c *gin.Context) {
	//	seelog.Info("Reload Sp Cache")
	//	spCache = make(map[string]Sp)
	//	spInitCache(g)
	//	utils.OK.OutJSON(c, nil)
	//})
	//if !gin.IsDebugging() {
	//	spInitCache(g)
	//}
	//if auth != nil {
	//	eng.POST(uri+"/*sp", func(c *gin.Context) {
	//		if auth(c) {
	//			spName := spPrefix + simpleNameReg.ReplaceAllString(c.Param("sp"), "")
	//			var sp Sp
	//			if gin.IsDebugging() {
	//				sp, _ = LoadSpFromDb(g, spName)
	//			} else {
	//				sp = spCache[spName]
	//			}
	//			data, err := GinSp(g, sp, c)
	//			if err == nil {
	//				utils.OK.OutJSON(c, data)
	//			} else {
	//				seelog.Error("数据存储过程错误:"+spName, "\n\t", err)
	//				utils.StErrorDb.OutJSON(c, nil)
	//			}
	//		} else {
	//			utils.StErrorToken.OutJSON(c, nil)
	//		}
	//	})
	//}
}
