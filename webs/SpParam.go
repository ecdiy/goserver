package webs

import (
	"github.com/cihub/seelog"
	"utils/http"
	"utils/xml"
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

/*--微信获取用户Id
1. db appId,secret.
2. http get openId
3. query openId to userId.
*/
func (ws *SpWeb) WxParam(ctx *Param, p *SpParam) (interface{}, int) {
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

func (ws *SpWeb) GinRpcParam(ele *xml.Element) {

}

func (ws *SpWeb) GinDbParam(ele *xml.Element) {
	ver := &RpcUser{Sql: ele.MustAttr("Sql"), Gpa: ws.Gpa}
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
		return 0, 401
	}
}
