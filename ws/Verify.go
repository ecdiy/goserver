package ws

import (
	"github.com/gin-gonic/gin"
	"strings"
	"context"
)

func verify(c *gin.Context, tokenName string) *Web {
	auth := &Web{}
	auth.Ua = GetUa(c)
	auth.Context = c
	auth.Out = make(map[string]interface{})
	sut, e := c.Cookie(auth.Ua + tokenName)
	if e == nil && len(sut) > 1 {
		idx := strings.Index(sut, "_")
		if idx > 0 {
			RpcUser(func(client RpcUserClient, ctx context.Context) {
				sc, _ := client.Verify(ctx, &Token{Token: sut, Ua: auth.Ua})
				if sc.Result {
					auth.UserId = sc.UserId
					auth.Username = sc.Username
					auth.Score = sc.Score
					auth.Auth = true
				} else {
					auth.Auth = false
				}
			})
		} else {
			auth.Auth = false
		}
	} else {
		auth.Auth = false
	}
	return auth
}

func VerifyRpc(c *gin.Context) *Web {
	return verify(c, "Token")
}

func VerifyAdmin(c *gin.Context) *Web {
	return verify(c, "Admin")
}
