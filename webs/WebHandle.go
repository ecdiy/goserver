package webs

import (
	"github.com/gin-gonic/gin"
	"context"
)

func Post(Gin *gin.Engine, url string, fun func(param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		wb := NewParam(c)
		fun(wb)
		c.JSON(200, wb.Out)
	})
}

func PostRpc(Gin *gin.Engine, rpc *RpcUser, tokenName, url string, fun func(ub *UserBase, param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		wb := NewParam(c)
		ub, _ := rpc.Verify(nil, &Token{Token: wb.String(tokenName), Ua: wb.Ua})
		fun(ub, wb)
		c.JSON(200, wb.Out)
	})
}
func PostHost(Gin *gin.Engine, RpcHost, tokenName, url string, fun func(ub *UserBase, param *Param)) {
	Gin.POST(url, func(c *gin.Context) {
		var ub *UserBase
		wb := NewParam(c)
		rpcUser(RpcHost, func(client RpcUserClient, ctx context.Context) {
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


//func GinGetAcao(Gin *gin.Engine, relativePath string, hand HandlerResult) {
//	Gin.GET(relativePath, func(c *gin.Context) {
//		res := hand(c)
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Header("Content-Type", "application/json; charset=utf-8")
//		c.JSON(http.StatusOK, res)
//	})
//}
