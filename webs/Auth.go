package webs

import (
	"encoding/json"
	"context"
	"strconv"
	"utils/gpa"
)

type RpcUser struct {
	Sql string
	Gpa *gpa.Gpa
}

func (s *RpcUser) Verify(c context.Context, in *Token) (*UserBase, error) {
	ub := &UserBase{}
	if len(in.Token) > 1 {
		v, b := TokenMap[in.Token]
		if b {
			ub.UserId = v
			setUb(ub)
		} else {
			m, b, ee := s.Gpa.QueryMapStringString(s.Sql, in.Token, in.Ua)
			if ee == nil && b {
				uId, _ := strconv.ParseInt(m["UserId"], 10, 0)
				TokenMap[in.Token] = uId
				ub.UserId = uId
				setUb(ub)
			} else {
				ub.Result = false
			}
		}
	}
	return ub, nil
}

func setUb(ub *UserBase) {
	ub.Result = true
	m, mb := UserMap[ub.UserId]
	if mb {
		un, unb := m["Username"]
		if unb {
			ub.Username = un
		}
		if len(m) > 2 {
			be, _ := json.Marshal(m)
			ub.AppendJson = string(be)
		}
	}
}

func GetAuthByRpc(rpc *RpcUser, tokenName string) func(param *Param) *UserBase {
	return func(param *Param) *UserBase {
		param.Context.Set("CallAuth", 1)
		tokenVal := param.String(tokenName)
		if len(tokenVal) > 1 {
			ub, _ := rpc.Verify(nil, &Token{Token: tokenVal, Ua: param.Ua})
			ubx(ub, param)
			return ub
		} else {
			return &UserBase{Result: false}
		}
	}
}
func ubx(ub *UserBase, param *Param) {
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
			ubx(ub, param)
			return ub
		} else {
			return &UserBase{Result: false}
		}
	}
}
