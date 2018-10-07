package webs

import (
	"encoding/json"
	"context"
)

func GetAuthByRpc(rpc *RpcUser, tokenName string) func(param *Param) *UserBase {
	return func(param *Param) *UserBase {
		param.Context.Set("CallAuth", 1)
		tokenVal := param.String(tokenName)
		if len(tokenVal) > 1 {
			ub, _ := rpc.Verify(nil, &Token{Token: tokenVal, Ua: param.Ua})
			if ub.Result {
				param.Context.Set("UserId", ub.UserId)
				param.Context.Set("Username", ub.Username)
				if len(ub.AppendJson) > 1 {
					var data map[string]interface{}
					je := json.Unmarshal([]byte(ub.AppendJson), &data)
					if je == nil {
						for k, v := range data {
							param.Context.Set(k, v)
						}
					}
				}
			}
			return ub
		} else {
			return &UserBase{Result: false}
		}
	}
}

func GetAuthFunByHost(RpcUserHost, tokenName string) func(param *Param) *UserBase {
	return func(param *Param) *UserBase {
		param.Context.Set("CallAuth", 1)
		tokenVal := param.String(tokenName)
		if len(tokenVal) > 1 {
			var ub *UserBase
			rpcUser(RpcUserHost, func(client RpcUserClient, ctx context.Context) {
				ub, _ = client.Verify(ctx, &Token{Token: tokenVal, Ua: param.Ua})
			})
			if ub.Result {
				param.Context.Set("UserId", ub.UserId)
				param.Context.Set("Username", ub.Username)
				if len(ub.AppendJson) > 1 {
					var data map[string]interface{}
					je := json.Unmarshal([]byte(ub.AppendJson), &data)
					if je == nil {
						for k, v := range data {
							param.Context.Set(k, v)
						}
					}
				}
			}
			return ub
		} else {
			return &UserBase{Result: false}
		}
	}
}
