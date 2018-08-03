package ws

import (
	"github.com/gin-gonic/gin"
	"strings"
	"context"
)

func VerifyRpc(c *gin.Context) *Web {
	auth := &Web{}
	auth.Ua = GetUa(c)
	auth.Context = c
	auth.Out = make(map[string]interface{})
	sut, e := c.Cookie(auth.Ua + "Token")
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

func VerifyAdmin(c *gin.Context) *Web {
	auth := &Web{}
	auth.Ua = GetUa(c)
	auth.Context = c
	auth.Out = make(map[string]interface{})
	sut, e := c.Cookie(auth.Ua + "Admin")
	if e == nil && len(sut) > 1 {
		idx := strings.Index(sut, "_")
		if idx > 0 {
			RpcAdmin(func(client RpcAdminClient, ctx context.Context) {
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
