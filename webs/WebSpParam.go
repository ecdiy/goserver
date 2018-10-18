package webs

import (
	"github.com/cihub/seelog"
	"utils/xml"
	"context"
)

//--
func ParamUa(wb *Param, p *SpParam) (interface{}, int) {
	return wb.Ua, 200
}

func ParamIn(ctx *Param, p *SpParam) (interface{}, int) {
	v := ctx.String(p.ParamName)
	if v == "" {
		v = p.DefaultVal
	}
	return v, 200
}

/*--微信获取用户Id
1. db appId,secret.
2. http get openId
3. query openId to userId.
*/
func (ws *SpWeb) ParamWx(ele *xml.Element, data map[string]interface{}) {
	//wb := NewParam(ctx.Context)
	//res := make(map[string]interface{})
	//c := SpCall("wx", wb, res, nil, true)
	//if c == 200 {
	//	wx := res["wx"].(map[string]string)
	//	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + wx["MpAppId"] +
	//		"&secret=" + wx["MpSecret"] +
	//		"&js_code=" + ctx.String("js_code") + "&grant_type=authorization_code"
	//	h := &http.Http{}
	//	_, e := h.GetMap(url)
	//	if e == nil {
	//
	//	}
	//}
}

func (ws *SpWeb) ParamGinRpc(ele *xml.Element, data map[string]interface{}) {

}

func (ws *SpWeb) ParamWk(ele *xml.Element, data map[string]interface{}) {
	RpcRef, rr := ele.AttrValue("RpcRef")
	var ver *RpcUser
	if rr {
		ver = data[RpcRef].(*RpcUser)
	} else {
		ver = &RpcUser{Sql: ele.MustAttr("Sql"), Gpa: ws.Gpa}
	}
	tkName := ele.MustAttr("TokenName")
	ws.SpParamDoMap[ele.MustAttr("Prefix")] = func(wb *Param, p *SpParam) (interface{}, int) {
		v, b := wb.Context.Get(p.ParamName)
		if b {
			return v, 200
		}
		ver.Verify(nil, &Token{Token: wb.String(tkName), Ua: wb.Ua})
		v2, b2 := wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		seelog.Error("ctx.Get not find.", p.ParamName)
		return p.DefaultVal, 200
	}
}

func (ws *SpWeb) ParamGin(ele *xml.Element, data map[string]interface{}) {

	tkName := ele.MustAttr("TokenName")
	prefix := ele.MustAttr("Prefix")

	RpcRef, rr := ele.AttrValue("RpcRef")
	var ver *RpcUser
	if rr {
		ver = data[RpcRef].(*RpcUser)
	} else {
		RpcHost, rhb := ele.AttrValue("RpcHost")
		if rhb {
			ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
				v, b := wb.Context.Get(p.ParamName)
				if b {
					return v, 200
				}
				var ub *UserBase
				rpcUser(RpcHost, func(client RpcUserClient, ctx context.Context) {
					ub, _ = client.Verify(ctx, &Token{Token: wb.String(tkName), Ua: wb.Ua})
				})
				if ub.Result {
					ubx(ub, wb)
					v2, b2 := wb.Context.Get(p.ParamName)
					if b2 {
						return v2, 200
					}
				}
				seelog.Error("ctx.Get not find.", p.ParamName)
				return 0, 401
			}
			return
		}
		ver = &RpcUser{Sql: ele.MustAttr("Sql"), Gpa: ws.Gpa}
	}
	ws.SpParamDoMap[prefix] = func(wb *Param, p *SpParam) (interface{}, int) {
		v, b := wb.Context.Get(p.ParamName)
		if b {
			return v, 200
		}
		ver.Verify(nil, &Token{Token: wb.String(tkName), Ua: wb.Ua})
		v2, b2 := wb.Context.Get(p.ParamName)
		if b2 {
			return v2, 200
		}
		seelog.Error("ctx.Get not find.", p.ParamName)
		return 0, 401
	}

}
