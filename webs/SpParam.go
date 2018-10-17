package webs

import (
	"github.com/cihub/seelog"
	"utils/http"
)

//--
func UaParam(wb *Param, p *SpParam) (interface{}, int) {
	return wb.Ua, 200
}

func InParam(ctx *Param, p *SpParam) (interface{}, int) {
	v := ctx.String(p.ParamName)
	if v == "" {
		v = p.DefaultVal
	}
	return v, 200
}

func (sp *Sp) GkParam(ctx *Param, p *SpParam) (interface{}, int) {
	_, CallAuth := ctx.Context.Get("CallAuth")
	if !CallAuth && sp.Auth != nil {
		sp.Auth(ctx)
	}
	v, b := ctx.Context.Get(p.ParamName)
	if b {
		return v, 200
	} else {
		seelog.Error("ctx.Get not find.", p.ParamName)
		return p.DefaultVal, 200
	}
}

func (sp *Sp) GinParam(ctx *Param, p *SpParam) (interface{}, int) {
	v, b := ctx.Context.Get(p.ParamName)
	if b {
		return v, 200
	} else {
		_, CallAuth := ctx.Context.Get("CallAuth")
		if !CallAuth && sp.Auth != nil {
			auth := sp.Auth(ctx)
			if !auth.Result {
				seelog.Error("Gin获取参数值出错：SpName=", sp.Name, ";ParamName=", p.ParamName)
				return "", 401
			}
			v, b := ctx.Context.Get(p.ParamName)
			if b {
				return v, 200
			} else {
				seelog.Error("获取Gin参数错误:", p.ParamName)
				return "", 401
			}
		}
		seelog.Error("ctx.Get not find.", p.ParamName)
		return "", 401
	}
}

/*--微信获取用户Id
1. db appId,secret.
2. http get openId
3. query openId to userId.
*/
func WxParam(ctx *Param, p *SpParam) (interface{}, int) {
	wb := NewParam(ctx.Context)
	res := make(map[string]interface{})
	c := SpCall("wx", wb, res, nil, true)
	if c == 200 {
		wx := res["wx"].(map[string]string)
		url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + wx["MpAppId"] +
			"&secret=" + wx["MpSecret"] +
			"&js_code=" + ctx.String("js_code") + "&grant_type=authorization_code"
		h := &http.Http{}
		_, e := h.GetMap(url)
		if e == nil {

		}
	}
	return 0, 401
}

//--
