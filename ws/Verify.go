package ws

import (
	"context"
	"strings"
)

func VerifyRpc(auth *Web) {
	sut, e := auth.Context.Cookie(auth.Ua + "Token")
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
}

func VerifyAdmin(auth *Web) {
	sut, e := auth.Context.Cookie(auth.Ua + "Admin")
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
}
