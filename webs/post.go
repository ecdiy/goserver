package webs

import (
	"github.com/gin-gonic/gin"
	"context"
)

//func Post(Gin *gin.Engine, url string, fun func(param *Param, res map[string]interface{})) {
//	Gin.POST(url, func(c *gin.Context) {
//		param := NewParam(c)
//
//		res := make(map[string]interface{})
//		fun(param, res)
//		res["now"] = time.Now().Format("2006-01-02T15:04:05Z")
//		c.JSON(200, res)
//
//	})
//}

func PostRpc(Gin *gin.Engine, rpc *RpcUser, tokenName, url string, fun func(ub *UserBase, param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		wb := NewParam(c)
		ub, _ := rpc.Verify(nil, &Token{Token: wb.String(tokenName), Ua: wb.Ua})
		fun(ub, wb)
		c.JSON(200, wb.Out)
	})
}
func PostHost(Gin *gin.Engine, RpcUserHost, tokenName, url string, fun func(ub *UserBase, param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		var ub *UserBase
		wb := NewParam(c)
		rpcUser(RpcUserHost, func(client RpcUserClient, ctx context.Context) {
			ub, _ = client.Verify(ctx, &Token{Token: wb.String(tokenName), Ua: wb.Ua})
		})
		fun(ub, wb)
		c.JSON(200, wb.Out)
	})
}
func AuthHost(Gin *gin.Engine, RpcUserHost, tokenName, url string, fun func(ub *UserBase, param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		rpcUser(RpcUserHost, func(client RpcUserClient, ctx context.Context) {
			wb := NewParam(c)
			ub, _ := client.Verify(ctx, &Token{Token: wb.String(tokenName), Ua: wb.Ua})
			if ub.Result {
				fun(ub, wb)
				c.JSON(200, wb.Out)
			} else {
				c.AbortWithStatus(401)
			}
		})
	})
}

func AuthRpc(Gin *gin.Engine, rpc *RpcUser, tokenName, url string, fun func(ub *UserBase, param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		wb := NewParam(c)
		ub, _ := rpc.Verify(nil, &Token{Token: wb.String(tokenName), Ua: wb.Ua})
		if ub.Result {
			fun(ub, wb)
			c.JSON(200, wb.Out)
		} else {
			c.AbortWithStatus(401)
		}
	})
}

//
//func Auth(Gin *gin.Engine, url string, fun func(userId int64, param *Param, res map[string]interface{}), verify func(c *gin.Context) (bool, int64)) {
//	Gin.POST(url, func(c *gin.Context) {
//		auth, userId := verify(c)
//		if auth {
//			param := NewParam(c)
//
//			param.Auth = c.Keys
//			res := make(map[string]interface{})
//			fun(userId, param, res)
//			res["now"] = time.Now().Format("2006-01-02T15:04:05Z")
//			c.JSON(200, res)
//
//		} else {
//			c.AbortWithStatus(401)
//		}
//	})
//}

//func GinGetAcao(Gin *gin.Engine, relativePath string, hand HandlerResult) {
//	Gin.GET(relativePath, func(c *gin.Context) {
//		res := hand(c)
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Header("Content-Type", "application/json; charset=utf-8")
//		c.JSON(http.StatusOK, res)
//	})
//}
